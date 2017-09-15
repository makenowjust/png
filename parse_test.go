package png

import (
	"testing"

	"strings"
)

func TestParse(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		p, err := Parse("")

		if err == nil {
			t.Fatalf("succeeded in Parse(): %+#v", p)
		}

		if msg := err.Error(); msg != "invalid URL: \"\" (empty)" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		p, err := Parse(":")

		if err == nil {
			t.Fatalf("succeeded in Parse(): %+#v", p)
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in parsing URL: \":\"") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Unknown Schme", func(t *testing.T) {
		p, err := Parse("invalid://")

		if err == nil {
			t.Fatalf("succeeded in Parse(): %+#v", p)
		}

		if msg := err.Error(); msg != "unknown scheme: invalid" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}

func TestParseHTTPURL(t *testing.T) {
	t.Run("No Scheme", func(t *testing.T) {
		p, err := Parse("localhost")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		hp, ok := p.(*HTTPPinger)
		if !ok {
			t.Fatalf("failed in casting to *HTTPPinger: %+#v", p)
		}

		if hp.url.String() != "http://localhost" {
			t.Fatalf("unexpected result: %#v", hp.url.String())
		}
	})

	t.Run("Opaque", func(t *testing.T) {
		p, err := Parse("localhost:8080")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		hp, ok := p.(*HTTPPinger)
		if !ok {
			t.Fatalf("failed in casting to *HTTPPinger: %+#v", p)
		}

		if hp.url.String() != "http://localhost:8080" {
			t.Fatalf("unexpected result: %#v", hp.url.String())
		}
	})

	t.Run("No Host", func(t *testing.T) {
		p, err := Parse("http://:8080")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		hp, ok := p.(*HTTPPinger)
		if !ok {
			t.Fatalf("failed in casting to *HTTPPinger: %+#v", p)
		}

		if hp.url.String() != "http://127.0.0.1:8080" {
			t.Fatalf("unexpected result: %#v", hp.url.String())
		}
	})

	t.Run("HTTP", func(t *testing.T) {
		p, err := Parse("http://localhost:8080")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		hp, ok := p.(*HTTPPinger)
		if !ok {
			t.Fatalf("failed in casting to *HTTPPinger: %+#v", p)
		}

		if hp.url.String() != "http://localhost:8080" {
			t.Fatalf("unexpected result: %#v", hp.url.String())
		}
	})

	t.Run("HTTPS", func(t *testing.T) {
		p, err := Parse("https://localhost:8080")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		hp, ok := p.(*HTTPPinger)
		if !ok {
			t.Fatalf("failed in casting to *HTTPPinger: %+#v", p)
		}

		if hp.url.String() != "https://localhost:8080" {
			t.Fatalf("unexpected result: %#v", hp.url.String())
		}
	})
}

func TestParseWebSocketURL(t *testing.T) {
	t.Run("WS", func(t *testing.T) {
		p, err := Parse("ws://localhost:8080")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		wp, ok := p.(*WebSocketPinger)
		if !ok {
			t.Fatalf("failed in casting to *WebSocketPinger: %+#v", p)
		}

		if wp.url.String() != "ws://localhost:8080" {
			t.Fatalf("unexpected result: %#v", wp.url.String())
		}
	})

	t.Run("WSS", func(t *testing.T) {
		p, err := Parse("wss://localhost:8080")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		wp, ok := p.(*WebSocketPinger)
		if !ok {
			t.Fatalf("failed in casting to *WebSocketPinger: %+#v", p)
		}

		if wp.url.String() != "wss://localhost:8080" {
			t.Fatalf("unexpected result: %#v", wp.url.String())
		}
	})
}

func TestParseTCPURL(t *testing.T) {
	t.Run("TCP", func(t *testing.T) {
		p, err := Parse("tcp://8.8.8.8:53")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		tp, ok := p.(*TCPPinger)
		if !ok {
			t.Fatalf("failed in casting to *TCPPinger: %+#v", p)
		}

		if tp.network != "tcp" || tp.addr != "8.8.8.8:53" {
			t.Fatalf("unexpected result: %+#v", tp)
		}
	})

	t.Run("TCP4", func(t *testing.T) {
		p, err := Parse("tcp4://8.8.8.8:53")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		tp, ok := p.(*TCPPinger)
		if !ok {
			t.Fatalf("failed in casting to *TCPPinger: %+#v", p)
		}

		if tp.network != "tcp4" || tp.addr != "8.8.8.8:53" {
			t.Fatalf("unexpected result: %+#v", tp)
		}
	})

	t.Run("TCP6", func(t *testing.T) {
		p, err := Parse("tcp6://[2001:4860:4860::8888]:53")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		tp, ok := p.(*TCPPinger)
		if !ok {
			t.Fatalf("failed in casting to *TCPPinger: %+#v", p)
		}

		if tp.network != "tcp6" || tp.addr != "[2001:4860:4860::8888]:53" {
			t.Fatalf("unexpected result: %+#v", tp)
		}
	})
}

func TestParseMySQLURL(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		p, err := Parse("mysql://root@localhost:13306/")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		mp, ok := p.(*MySQLPinger)
		if !ok {
			t.Fatalf("failed in casting to *MySQLPinger: %+#v", p)
		}

		if mp.url.String() != "mysql://root@localhost:13306/" {
			t.Fatalf("unexpected result: %#v", mp.url.String())
		}
	})

	t.Run("No Host", func(t *testing.T) {
		p, err := Parse("mysql://root@/")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		mp, ok := p.(*MySQLPinger)
		if !ok {
			t.Fatalf("failed in casting to *MySQLPinger: %+#v", p)
		}

		if mp.url.String() != "mysql://root@127.0.0.1:3306/" {
			t.Fatalf("unexpected result: %#v", mp.url.String())
		}
	})
}

func TestParsePostgresURL(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		p, err := Parse("postgres://root@localhost:15043/pg?sslmode=require")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		pp, ok := p.(*PostgresPinger)
		if !ok {
			t.Fatalf("failed in casting to *PostgresPinger: %+#v", p)
		}

		if pp.url.String() != "postgres://root@localhost:15043/pg?sslmode=require" {
			t.Fatalf("unexpected result: %#v", pp.url.String())
		}
	})

	t.Run("No SSLMode", func(t *testing.T) {
		p, err := Parse("postgres://root@localhost:15043/pg")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		pp, ok := p.(*PostgresPinger)
		if !ok {
			t.Fatalf("failed in casting to *PostgresPinger: %+#v", p)
		}

		if pp.url.String() != "postgres://root@localhost:15043/pg?sslmode=disable" {
			t.Fatalf("unexpected result: %#v", pp.url.String())
		}
	})

	t.Run("No Table", func(t *testing.T) {
		p, err := Parse("postgres://root@localhost:15043/")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		pp, ok := p.(*PostgresPinger)
		if !ok {
			t.Fatalf("failed in casting to *PostgresPinger: %+#v", p)
		}

		if pp.url.String() != "postgres://root@localhost:15043/postgres?sslmode=disable" {
			t.Fatalf("unexpected result: %#v", pp.url.String())
		}
	})
}

func TestParseRedisURL(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		p, err := Parse("redis://:password@localhost:16379/42")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		rp, ok := p.(*RedisPinger)
		if !ok {
			t.Fatalf("failed in casting to *RedisPinger: %+#v", p)
		}

		if rp.addr != "localhost:16379" {
			t.Fatalf("unexpected addr: %#v", rp.addr)
		}

		if rp.password != "password" {
			t.Fatalf("unexpected password: %#v", rp.password)
		}

		if rp.db != 42 {
			t.Fatalf("unexpected db number: %v", rp.db)
		}
	})

	t.Run("Invalid DB", func(t *testing.T) {
		p, err := Parse("redis://localhost/invalid")

		if err == nil {
			t.Fatal("succeeded in Parse()", p)
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "invalid db number: \"invalid\": ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}

func TestParseAMQPURL(t *testing.T) {
	p, err := Parse("amqp://localhost:5672")

	if err != nil {
		t.Fatalf("failed in Parse(): %+#v", err)
	}

	ap, ok := p.(*AMQPPinger)
	if !ok {
		t.Fatalf("failed in casting to *WebSocketPinger: %+#v", p)
	}

	if ap.url.String() != "amqp://localhost:5672" {
		t.Fatalf("unexpected result: %#v", ap.url.String())
	}
}
