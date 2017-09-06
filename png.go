package png

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Timeout struct {
	Err error
}

func (timeout *Timeout) Error() string {
	return timeout.Err.Error()
}

type Pinger interface {
	Addr() (hostname string, port int, err error)
	Ping(ctx context.Context) error
}

func PingWithTimeout(p Pinger, timeout time.Duration) (elapsed time.Duration, err error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error)
	go func() {
		done <- p.Ping(ctx)
	}()

	select {
	case <-ctx.Done():
		elapsed = time.Since(start)
		err = &Timeout{Err: ctx.Err()}
	case err = <-done:
		elapsed = time.Since(start)
	}

	return
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
			fallthrough
		case "ws":
			port = 80
		case "https":
			fallthrough
		case "wss":
			port = 443
		case "postgres":
			port = 5432
		case "amqp":
			port = 5672
		default:
			err = errors.Errorf("invalid scheme: %s", p.url.Scheme)
		}
	}

	return
}
