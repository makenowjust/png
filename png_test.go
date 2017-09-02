package png

import (
	"testing"

	"strings"
)

func TestNewPinger(t *testing.T) {
	t.Run("No Schema", func(t *testing.T) {
		t.Run("No Port", func(t *testing.T) {
			p, err := NewPinger("127.0.0.1")
			if err != nil {
				t.Fatal("NewPinger failed", err)
			}

			hp, ok := p.(*HTTPPinger)
			if !ok {
				t.Fatal("failed to cast to HTTPPinger", hp)
			}

			if hp.URL != "http://127.0.0.1" {
				t.Fatalf("unexpected URL: %#v", hp.URL)
			}
		})

		t.Run("Port", func(t *testing.T) {
			p, err := NewPinger("localhost:8080")
			if err != nil {
				t.Fatal("NewPinger failed", err)
			}

			hp, ok := p.(*HTTPPinger)
			if !ok {
				t.Fatal("failed to cast to HTTPPinger", hp)
			}

			if hp.URL != "http://localhost:8080" {
				t.Fatalf("unexpected URL: %#v", hp.URL)
			}
		})
	})

	t.Run("HTTP", func(t *testing.T) {
		p, err := NewPinger("http://localhost:8080")
		if err != nil {
			t.Fatal("NewPinger failed", err)
		}

		hp, ok := p.(*HTTPPinger)
		if !ok {
			t.Fatal("failed to cast to HTTPPinger", hp)
		}

		if hp.URL != "http://localhost:8080" {
			t.Fatalf("unexpected URL: %+#v", hp.URL)
		}
	})

	t.Run("Redis", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			p, err := NewPinger("redis://:password@localhost:8080/42")
			if err != nil {
				t.Fatal("NewPinger failed", err)
			}

			rp, ok := p.(*RedisPinger)
			if !ok {
				t.Fatal("failed to cast to RedisPinger", rp)
			}

			if rp.Addr != "localhost:8080" || rp.Password != "password" || rp.DB != 42 {
				t.Fatalf("unexpected RedisPinger: %+#v", rp)
			}
		})

		t.Run("Invalid DB", func(t *testing.T) {
			_, err := NewPinger("redis://:password@localhost:8080/invalid_db")

			if err == nil {
				t.Fatal("succeed to create a pinger")
			}

			if msg := err.Error(); !strings.HasPrefix(msg, "invalid redis db: ") {
				t.Fatalf("unexpected error message: %#v", msg)
			}
		})
	})

	t.Run("Default Redis", func(t *testing.T) {
		p, err := NewPinger("redis")
		if err != nil {
			t.Fatal("NewPinger failed", err)
		}

		rp, ok := p.(*RedisPinger)
		if !ok {
			t.Fatal("failed to cast to RedisPinger", rp)
		}

		if rp.Addr != "localhost:6379" || rp.Password != "" || rp.DB != 0 {
			t.Fatalf("unexpected RedisPinger: %#v", rp)
		}
	})

	t.Run("Unsupported Scheme", func(t *testing.T) {
		_, err := NewPinger("unsupport://localhost")

		if err == nil {
			t.Fatal("succeed to create a pinger")
		}

		if msg := err.Error(); msg != "unsupported scheme: unsupport" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		_, err := NewPinger("::")

		if err == nil {
			t.Fatal("succeed to create a pinger")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "invalid URL: ::: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
