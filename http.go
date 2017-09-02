package png

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type HTTPPinger struct {
	URL string
}

func (pinger *HTTPPinger) Ping(ctx context.Context) error {
	req, err := http.NewRequest("HEAD", pinger.URL, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create HTTP request")
	}

	req.Header.Add("User-Agent", "png/0.0.0-dev")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed HTTP request")
	}

	if resp.StatusCode >= 400 {
		return errors.Errorf("failed HTTP request with %s", resp.Status)
	}

	return nil
}
