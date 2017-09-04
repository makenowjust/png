package png

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type TCPPinger struct {
	network  string
	hostname string
	port     int
}

func (p *TCPPinger) Addr() (string, int, error) {
	return p.hostname, p.port, nil
}

func (p *TCPPinger) Ping(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", p.hostname, p.port)
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, p.network, address)
	if err != nil {
		return errors.Wrapf(err, "failed in connecting %s on %s", address, p.network)
	}
	defer conn.Close()

	return nil
}
