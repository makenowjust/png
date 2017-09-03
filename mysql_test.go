package png

import (
	"testing"

	"net/url"
)

func TestMySQLURLToDSN(t *testing.T) {
	u, err := url.Parse("mysql://localhost")
	if err != nil {
		panic(err)
	}

	p := &MySQLPinger{urlPinger: &urlPinger{url: u}}

	if dsn := p.urlToDSN(); dsn != "tcp(localhost)" {
		t.Fatalf("unexpected dsn: %#v", dsn)
	}

	if dsn := p.urlToDSN(); dsn != "tcp(localhost)" {
		t.Fatalf("unexpected dsn: %#v", dsn)
	}
}
