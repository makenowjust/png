package png

import (
	"context"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func NewPinger(target string) (Pinger, error) {
	if target == "redis" {
		return &RedisPinger{Addr: "localhost:6379"}, nil
	}

	u, err := url.Parse(target)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid URL: %s", target)
	}

	switch u.Scheme {
	case "http":
		fallthrough
	case "https":
		return &HTTPPinger{URL: target}, nil

	case "redis":
		var password string
		if u.User != nil {
			if p, ok := u.User.Password(); ok {
				password = p
			}
		}

		var db int64
		if len(u.Path) > 2 {
			db, err = strconv.ParseInt(u.Path[1:], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "invalid redis db")
			}
		}

		return &RedisPinger{
			Addr:     u.Host,
			Password: password,
			DB:       int(db),
		}, nil

	case "":
		return &HTTPPinger{URL: "http://" + target}, nil

	default:
		// For "localhost:8080" case
		if u.Opaque != "" && "0" <= u.Opaque[0:1] && u.Opaque[0:1] <= "9" {
			return &HTTPPinger{URL: "http://" + target}, nil
		}

		return nil, errors.Errorf("unsupported scheme: %s", u.Scheme)
	}
}
