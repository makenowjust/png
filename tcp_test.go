package png

import (
	"testing"

	"context"
	"net"
	"strings"
)

func runTCPServer(network, addr string, handler func(net.Conn)) (*net.TCPListener, string) {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		panic(err)
	}

	s, err := net.ListenTCP("tcp", tcpAddr)
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

	return s, s.Addr().String()
}

func TestTCPPinger(t *testing.T) {
	t.Run("TCP", func(t *testing.T) {
		s, addr := runTCPServer("tcp", "localhost:0", func(conn net.Conn) {
			conn.Close()
		})
		defer s.Close()

		p := &TCPPinger{network: "tcp", addr: addr}
		err := p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("TCP4", func(t *testing.T) {
		s, addr := runTCPServer("tcp4", "localhost:0", func(conn net.Conn) {
			conn.Close()
		})
		defer s.Close()

		p := &TCPPinger{network: "tcp4", addr: addr}
		err := p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("TCP6", func(t *testing.T) {
		s, addr := runTCPServer("tcp6", "localhost:0", func(conn net.Conn) {
			conn.Close()
		})
		defer s.Close()

		p := &TCPPinger{network: "tcp6", addr: addr}
		err := p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("Unknown Network", func(t *testing.T) {
		p := &TCPPinger{network: "invalid", addr: "8.8.8.8:53"}
		err := p.Ping(context.Background())

		if err == nil {
			t.Fatal("succeeded in p.Ping()")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in connecting 8.8.8.8:53 on invalid") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
