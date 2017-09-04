package png

import (
	"testing"
)

func TestTCPPingerAddr(t *testing.T) {
	p := &TCPPinger{network: "tcp", hostname: "8.8.8.8", port: 53}
	hostname, port, err := p.Addr()

	if err != nil {
		t.Fatalf("failed in p.Addr(): %+#v", err)
	}

	if hostname != "8.8.8.8" {
		t.Fatalf("unexpected hostname: %#v", hostname)
	}

	if port != 53 {
		t.Fatalf("unexpected port: %v", port)
	}
}

func TestTCPPingerAddress(t *testing.T) {
	t.Run("IPv4", func(t *testing.T) {
		p := &TCPPinger{network: "tcp4", hostname: "8.8.8.8", port: 53}
		address := p.address()

		if address != "8.8.8.8:53" {
			t.Fatalf("unexpected address: %#v", address)
		}
	})

	t.Run("IPv6", func(t *testing.T) {
		p := &TCPPinger{network: "tcp6", hostname: "2001:4860:4860::8888", port: 53}
		address := p.address()

		if address != "[2001:4860:4860::8888]:53" {
			t.Fatalf("unexpected address: %#v", address)
		}
	})
}
