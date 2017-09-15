package png

import (
	"testing"

	"context"
	"net"
	"strings"
	"time"

	"github.com/alicebob/miniredis"
)

func runMiniredis() (*miniredis.Miniredis, string) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s, s.Addr()
}

func TestRedisPingerPing(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s, addr := runMiniredis()
		defer s.Close()

		p := &RedisPinger{addr: addr}
		err := p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("Password", func(t *testing.T) {
		s, addr := runMiniredis()
		defer s.Close()

		s.RequireAuth("password")

		p := &RedisPinger{addr: addr, password: "password"}
		err := p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		s, addr := runMiniredis()
		defer s.Close()

		p := &RedisPinger{addr: addr}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := p.Ping(ctx)

		if err == nil {
			t.Fatal("succeeded to ping")
		}

		if msg := err.Error(); msg != "failed in PING command: context canceled" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		s, addr := runTCPServer("tcp", "localhost:0", func(conn net.Conn) {
			time.Sleep(200 * time.Millisecond)
		})
		defer s.Close()

		p := &RedisPinger{addr: addr}
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		err := p.Ping(ctx)

		if err == nil {
			t.Fatal("succeeded to ping")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in PING command: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Pong", func(t *testing.T) {
		s, addr := runTCPServer("tcp", "localhost:0", func(conn net.Conn) {
			conn.Write([]byte("+BANG\r\n"))
		})
		defer s.Close()

		p := &RedisPinger{addr: addr}
		err := p.Ping(context.Background())

		if err == nil {
			t.Fatal("succeeded to ping")
		}

		if msg := err.Error(); msg != "invalid redis response: \"BANG\"" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
