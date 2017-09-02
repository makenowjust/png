package png

import (
	"testing"

	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func TestHTTPPinger(t *testing.T) {
	t.Run("Status", func(t *testing.T) {
		for _, success := range []int{
			// 1xx
			// http.StatusContinue,
			http.StatusSwitchingProtocols,
			http.StatusProcessing,

			// 2xx
			http.StatusOK,
			http.StatusCreated,
			http.StatusAccepted,
			http.StatusNonAuthoritativeInfo,
			http.StatusNoContent,
			http.StatusResetContent,
			http.StatusPartialContent,
			http.StatusMultiStatus,
			http.StatusAlreadyReported,
			http.StatusIMUsed,

			// 3xx
			http.StatusMultipleChoices,
			// http.StatusMovedPermanently,
			// http.StatusFound,
			// http.StatusSeeOther,
			http.StatusNotModified,
			http.StatusUseProxy,

			http.StatusTemporaryRedirect,
			http.StatusPermanentRedirect,
		} {
			t.Run(fmt.Sprintf("%d %s", success, http.StatusText(success)), func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(success)
				})
				s := httptest.NewServer(handler)
				defer s.Close()

				p := &HTTPPinger{URL: s.URL}
				err := p.Ping(context.Background())
				if err != nil {
					t.Fatalf("failed to ping: %+#v", err)
				}
			})
		}

		for _, fail := range []int{
			// 4xx
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusPaymentRequired,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusMethodNotAllowed,
			http.StatusNotAcceptable,
			http.StatusProxyAuthRequired,
			http.StatusRequestTimeout,
			http.StatusConflict,
			http.StatusGone,
			http.StatusLengthRequired,
			http.StatusPreconditionFailed,
			http.StatusRequestEntityTooLarge,
			http.StatusRequestURITooLong,
			http.StatusUnsupportedMediaType,
			http.StatusRequestedRangeNotSatisfiable,
			http.StatusExpectationFailed,
			http.StatusTeapot,
			http.StatusUnprocessableEntity,
			http.StatusLocked,
			http.StatusFailedDependency,
			http.StatusUpgradeRequired,
			http.StatusPreconditionRequired,
			http.StatusTooManyRequests,
			http.StatusRequestHeaderFieldsTooLarge,
			http.StatusUnavailableForLegalReasons,

			// 5xx
			http.StatusNotImplemented,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
			http.StatusHTTPVersionNotSupported,
			http.StatusVariantAlsoNegotiates,
			http.StatusInsufficientStorage,
			http.StatusLoopDetected,
			http.StatusNotExtended,
			http.StatusNetworkAuthenticationRequired,
		} {
			t.Run(fmt.Sprintf("%d %s", fail, http.StatusText(fail)), func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(fail)
				})
				s := httptest.NewServer(handler)
				defer s.Close()

				p := &HTTPPinger{URL: s.URL}
				err := p.Ping(context.Background())
				if err == nil {
					t.Fatal("succeed to ping")
				}

				if msg := err.Error(); msg != fmt.Sprintf("failed HTTP request with %d %s", fail, http.StatusText(fail)) {
					t.Fatalf("unexpected error message: %#v", msg)
				}
			})
		}
	})

	t.Run("Context", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Second)
			w.WriteHeader(http.StatusOK)
		})
		s := httptest.NewServer(handler)
		defer s.Close()

		p := &HTTPPinger{URL: s.URL}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := p.Ping(ctx)
		<-ctx.Done()

		if msg := err.Error(); msg != "failed HTTP request: Head "+s.URL+": context canceled" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
