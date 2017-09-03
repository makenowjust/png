package png

import (
	"testing"

	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/alicebob/miniredis"
)

func TestRedisPingerAddr(t *testing.T) {
	p := &RedisPinger{hostname: "localhost", port: 6379}
	hostname, port, err := p.Addr()

	if err != nil {
		t.Fatalf("failed in p.Addr(): %+#v", err)
	}

	if hostname != "localhost" {
		t.Fatalf("unexpected hostname: %#v", hostname)
	}

	if port != 6379 {
		t.Fatalf("unexpected port: %v", port)
	}
}

func runMiniredis() (*miniredis.Miniredis, string, int) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(s.Port())
	if err != nil {
		s.Close()
		panic(err)
	}

	return s, s.Host(), port
}

func runTCPServer(handler func(net.Conn)) (*net.TCPListener, string, int) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	s, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := s.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				handler(conn)
			}()
		}
	}()

	addr = s.Addr().(*net.TCPAddr)

	return s, addr.IP.String(), addr.Port
}

func TestRedisPingerPing(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s, hostname, port := runMiniredis()
		defer s.Close()

		p := &RedisPinger{hostname: hostname, port: port}
		err := p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("Password", func(t *testing.T) {
		s, hostname, port := runMiniredis()
		defer s.Close()

		s.RequireAuth("password")

		p := &RedisPinger{hostname: hostname, port: port, password: "password"}
		err := p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		s, hostname, port := runMiniredis()
		defer s.Close()

		p := &RedisPinger{hostname: hostname, port: port}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := p.Ping(ctx)

		if err == nil {
			t.Fatal("succeed to ping")
		}

		if msg := err.Error(); msg != "failed in PING command: context canceled" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		s, hostname, port := runTCPServer(func(conn net.Conn) {
			time.Sleep(200 * time.Millisecond)
		})
		defer s.Close()

		p := &RedisPinger{hostname: hostname, port: port}
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		err := p.Ping(ctx)

		if err == nil {
			t.Fatal("succeed to ping")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in PING command: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Pong", func(t *testing.T) {
		s, hostname, port := runTCPServer(func(conn net.Conn) {
			conn.Write([]byte("+BANG\r\n"))
		})
		defer s.Close()

		p := &RedisPinger{hostname: hostname, port: port}
		err := p.Ping(context.Background())

		if err == nil {
			t.Fatal("succeed to ping")
		}

		if msg := err.Error(); msg != "invalid redis response: \"BANG\"" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
