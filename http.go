package png

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type HTTPPinger struct {
	url *url.URL
}

func (p *HTTPPinger) Addr() (hostname string, port int, err error) {
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
		default:
			err = errors.Errorf("invalid scheme: %s", p.url.Scheme)
		}
	}

	return
}

func (p *HTTPPinger) Ping(ctx context.Context) error {
	req, err := http.NewRequest("HEAD", p.url.String(), nil)
	if err != nil {
		return errors.Wrap(err, "failed in creating HTTP request")
	}

	req.Header.Add("User-Agent", "png/0.0.0-dev")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed in HTTP request")
	}

	if resp.StatusCode >= 400 {
		return errors.Errorf("failed in HTTP request by %s", resp.Status)
	}

	return nil
}
