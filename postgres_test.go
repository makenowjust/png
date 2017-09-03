package png

import (
	"testing"

	"net/url"
)

func TestPostgresPingerAddr(t *testing.T) {
	t.Run("Port", func(t *testing.T) {
		u, err := url.Parse("postgres://localhost:5433")
		if err != nil {
			panic(err)
		}

		p := &PostgresPinger{urlPinger: &urlPinger{url: u}}
		hostname, port, err := p.Addr()

		if err != nil {
			t.Fatalf("failed in p.Addr(): %+#v", err)
		}

		if hostname != "localhost" {
			t.Fatalf("unexpected hostname: %#v", port)
		}

		if port != 5433 {
			t.Fatalf("unexpected port: %v", port)
		}
	})

	t.Run("No Port", func(t *testing.T) {
		u, err := url.Parse("postgres://localhost")
		if err != nil {
			panic(err)
		}

		p := &PostgresPinger{urlPinger: &urlPinger{url: u}}
		hostname, port, err := p.Addr()

		if err != nil {
			t.Fatalf("failed in p.Addr(): %+#v", err)
		}

		if hostname != "localhost" {
			t.Fatalf("unexpected hostname: %#v", port)
		}

		if port != 5432 {
			t.Fatalf("unexpected port: %v", port)
		}
	})
}
