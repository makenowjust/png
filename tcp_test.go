package png

import (
	"testing"
)

func TestTCPPingerAddr(t *testing.T) {
	p := &TCPPinger{hostname: "8.8.8.8", port: 53}
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
