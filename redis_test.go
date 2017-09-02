package png

import (
	"testing"

	"context"
	"net"
	"strings"
	"time"

	"github.com/alicebob/miniredis"
)

func TestRedisPinger(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		defer s.Close()

		p := &RedisPinger{Addr: s.Addr()}
		err = p.Ping(context.Background())
		if err != nil {
			t.Fatal("failed to ping: %+#v", err)
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		defer s.Close()

		p := &RedisPinger{Addr: s.Addr()}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err = p.Ping(ctx)

		if err == nil {
			t.Fatal("succeed to ping")
		}

		if msg := err.Error(); msg != "failed to ping to redis: context canceled" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:55555")
		if err != nil {
			panic(err)
		}
		s, err := net.ListenTCP("tcp", addr)
		if err != nil {
			panic(err)
		}
		defer s.Close()

		p := &RedisPinger{Addr: addr.String()}
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		err = p.Ping(ctx)

		if err == nil {
			t.Fatal("succeed to ping")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed to ping to redis: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Pong", func(t *testing.T) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:55555")
		if err != nil {
			panic(err)
		}
		s, err := net.ListenTCP("tcp", addr)
		if err != nil {
			panic(err)
		}
		go func() {
			conn, err := s.Accept()
			if err != nil {
				panic(err)
			}
			time.Sleep(100 * time.Millisecond)
			conn.Write([]byte("+BANG\r\n"))
		}()
		defer s.Close()

		p := &RedisPinger{Addr: addr.String()}
		err = p.Ping(context.Background())

		if err == nil {
			t.Fatal("succeed to ping")
		}

		if msg := err.Error(); msg != "invalid redis response: \"BANG\"" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
