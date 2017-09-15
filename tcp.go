package png

import (
	"context"
	"net"

	"github.com/pkg/errors"
)

type TCPPinger struct {
	network string
	addr    string
}

func (p *TCPPinger) Ping(ctx context.Context) error {
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, p.network, p.addr)
	if err != nil {
		return errors.Wrapf(err, "failed in connecting %s on %s", p.addr, p.network)
	}
	defer conn.Close()

	return nil
}
