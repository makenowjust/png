package png

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func NewPinger(target string) (Pinger, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid URL: %s", target)
	}

	switch u.Scheme {
	default:
		return nil, errors.Errorf("unsupported schema: %s", u.Scheme)
	}
}
