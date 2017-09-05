package png

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type HTTPPinger struct {
	*urlPinger
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
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.Errorf("failed in HTTP request by %s", resp.Status)
	}

	return nil
}
