package png

import (
	"testing"

	"strings"
)

func TestParse(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		p, err := Parse("")

		if err == nil {
			t.Fatalf("succeed in Parse(): %+#v", p)
		}

		if msg := err.Error(); msg != "cannot create pinger from empty string" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		p, err := Parse(":")

		if err == nil {
			t.Fatalf("succeed in Parse(): %+#v", p)
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in parsing URL: \":\"") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Invalid Schme", func(t *testing.T) {
		p, err := Parse("invalid://")

		if err == nil {
			t.Fatalf("succeed in Parse(): %+#v", p)
		}

		if msg := err.Error(); msg != "invalid scheme: invalid" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}

func TestParseForHTTP(t *testing.T) {
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

func TestParseForMySQL(t *testing.T) {
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

func TestParseForPostgres(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		p, err := Parse("postgres://root@localhost:15043/pg?sslmode=require")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		mp, ok := p.(*PostgresPinger)
		if !ok {
			t.Fatalf("failed in casting to *MySQLPinger: %+#v", p)
		}

		if mp.url.String() != "postgres://root@localhost:15043/pg?sslmode=require" {
			t.Fatalf("unexpected result: %#v", mp.url.String())
		}
	})

	t.Run("No SSLMode", func(t *testing.T) {
		p, err := Parse("postgres://root@localhost:15043/pg")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		mp, ok := p.(*PostgresPinger)
		if !ok {
			t.Fatalf("failed in casting to *MySQLPinger: %+#v", p)
		}

		if mp.url.String() != "postgres://root@localhost:15043/pg?sslmode=disable" {
			t.Fatalf("unexpected result: %#v", mp.url.String())
		}
	})

	t.Run("No Table", func(t *testing.T) {
		p, err := Parse("postgres://root@localhost:15043/")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		mp, ok := p.(*PostgresPinger)
		if !ok {
			t.Fatalf("failed in casting to *MySQLPinger: %+#v", p)
		}

		if mp.url.String() != "postgres://root@localhost:15043/postgres?sslmode=disable" {
			t.Fatalf("unexpected result: %#v", mp.url.String())
		}
	})
}

func TestParseForRedis(t *testing.T) {
	t.Run("No Port", func(t *testing.T) {
		p, err := Parse("redis://localhost")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		rp, ok := p.(*RedisPinger)
		if !ok {
			t.Fatalf("failed in casting to *RedisPinger: %+#v", p)
		}

		if rp.port != 6379 {
			t.Fatalf("unexpected port number: %v", rp.port)
		}
	})

	t.Run("All", func(t *testing.T) {
		p, err := Parse("redis://:password@localhost:16379/42")

		if err != nil {
			t.Fatalf("failed in Parse(): %+#v", err)
		}

		rp, ok := p.(*RedisPinger)
		if !ok {
			t.Fatalf("failed in casting to *RedisPinger: %+#v", p)
		}

		if rp.hostname != "localhost" {
			t.Fatalf("unexpected hostname: %#v", rp.hostname)
		}

		if rp.port != 16379 {
			t.Fatalf("unexpected port number: %v", rp.port)
		}

		if rp.password != "password" {
			t.Fatalf("unexpected password: %#v", rp.password)
		}

		if rp.db != 42 {
			t.Fatalf("unexpected db number: %v", rp.db)
		}
	})

	t.Run("Invalid Port", func(t *testing.T) {
		p, err := Parse("redis://localhost:invalid")

		if err == nil {
			t.Fatal("succeed in Parse()", p)
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "invalid port name: invalid: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Invalid DB", func(t *testing.T) {
		p, err := Parse("redis://localhost/invalid")

		if err == nil {
			t.Fatal("succeed in Parse()", p)
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "invalid db number: \"invalid\": ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
