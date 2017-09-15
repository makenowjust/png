package png

import (
	"context"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type WebSocketPinger struct {
	url *url.URL
}

func (ws *WebSocketPinger) Ping(ctx context.Context) error {
	done := make(chan error)

	go func() {
		dialer := &websocket.Dialer{
			NetDial: func(network, addr string) (net.Conn, error) {
				dialer := &net.Dialer{}
				if conn, err := dialer.DialContext(ctx, network, addr); err != nil {
					return nil, err
				} else {
					if t, ok := ctx.Deadline(); ok {
						conn.SetDeadline(t)
						conn.SetReadDeadline(t)
						conn.SetWriteDeadline(t)
					}
					return conn, err
				}
			},
		}

		header := http.Header{}
		header.Add("User-Agent", "png/0.0.0-dev")
		conn, _, err := dialer.Dial(ws.url.String(), header)
		if err != nil {
			done <- errors.Wrap(err, "failed in opening WebSocket connection")
			return
		}
		defer conn.Close()

		// TODO: should does it check resp fileds?

		done <- nil
	}()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "failed in WebSocket ping")
	case err := <-done:
		return err
	}
}
