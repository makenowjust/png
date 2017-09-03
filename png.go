package png

import (
	"context"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type Pinger interface {
	Addr() (hostname string, port int, err error)
	Ping(ctx context.Context) error
}

type urlPinger struct {
	url *url.URL
}

func (p *urlPinger) Addr() (hostname string, port int, err error) {
	hostname = p.url.Hostname()

	if portString := p.url.Port(); portString != "" {
		port, err = strconv.Atoi(portString)
		err = errors.Wrap(err, "failed in parsing port number")
	} else {
		switch p.url.Scheme {
		case "http":
			port = 80
		case "https":
			port = 443
		case "postgres":
			port = 5432
		default:
			err = errors.Errorf("invalid scheme: %s", p.url.Scheme)
		}
	}

	return
}
