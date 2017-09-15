package png

import (
	"testing"

	"context"
	"net/http"
	"net/url"
	"time"
	"strings"

	"github.com/gorilla/websocket"
)

func TestWebSocketPinger(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s, u := runHTTPServer(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			conn, err := upgrader.Upgrade(w, r, r.Header)
			if err != nil {
				t.Fatalf("failed in upgrading to WebSocket: %+#v", err)
			}
			defer conn.Close()
		})
		defer s.Close()

		u.Scheme = "ws"
		p := &WebSocketPinger{url: u}
		err := p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("Timeoout", func(t *testing.T) {
		s, u := runHTTPServer(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
		})
		defer s.Close()

		u.Scheme = "ws"
		p := &WebSocketPinger{url: u}
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		err := p.Ping(ctx)
		if err == nil {
			t.Fatal("succeeded in p.Ping()")
		}

    if msg := err.Error(); msg != "failed in WebSocket ping: context deadline exceeded" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		u, err := url.Parse("ws://not_found:not_found")
		if err != nil {
			panic(err)
		}

		p := &WebSocketPinger{url: u}
		err = p.Ping(context.Background())
		if err == nil {
			t.Fatal("succeeded in p.Ping()")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in opening WebSocket connection: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
