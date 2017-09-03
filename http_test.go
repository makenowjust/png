package png

import (
	"testing"

	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func TestHTTPPingerAddr(t *testing.T) {
	t.Run("Port", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			u, err := url.Parse("http://localhost:8080")
			if err != nil {
				panic(err)
			}

			p := &HTTPPinger{url: u}
			hostname, port, err := p.Addr()

			if err != nil {
				t.Fatalf("failed in p.Addr(): %+#v", err)
			}

			if hostname != "localhost" {
				t.Fatalf("unexpected hostname: %#v", port)
			}

			if port != 8080 {
				t.Fatalf("unexpected port: %v", port)
			}
		})

		t.Run("Invalid", func(t *testing.T) {
			u, err := url.Parse("http://localhost:invalid")
			if err != nil {
				panic(err)
			}

			p := &HTTPPinger{url: u}
			hostname, port, err := p.Addr()

			if err == nil {
				t.Fatal("succeed in p.Addr()", hostname, port)
			}

			if msg := err.Error(); !strings.HasPrefix(msg, "failed in parsing port number: ") {
				t.Fatalf("unexpected error message: %#v", msg)
			}
		})
	})

	t.Run("Wellknown", func(t *testing.T) {
		t.Run("HTTP", func(t *testing.T) {
			u, err := url.Parse("http://localhost")
			if err != nil {
				panic(err)
			}

			p := &HTTPPinger{url: u}
			hostname, port, err := p.Addr()

			if err != nil {
				t.Fatalf("failed in p.Addr(): %+#v", err)
			}

			if hostname != "localhost" {
				t.Fatalf("unexpected hostname: %#v", port)
			}

			if port != 80 {
				t.Fatalf("unexpected port: %v", port)
			}
		})

		t.Run("HTTPS", func(t *testing.T) {
			u, err := url.Parse("https://localhost")
			if err != nil {
				panic(err)
			}

			p := &HTTPPinger{url: u}
			hostname, port, err := p.Addr()

			if err != nil {
				t.Fatalf("failed in p.Addr(): %+#v", err)
			}

			if hostname != "localhost" {
				t.Fatalf("unexpected hostname: %#v", port)
			}

			if port != 443 {
				t.Fatalf("unexpected port: %v", port)
			}
		})

		t.Run("Invalid", func(t *testing.T) {
			u, err := url.Parse("invalid://localhost")
			if err != nil {
				panic(err)
			}

			p := &HTTPPinger{url: u}
			hostname, port, err := p.Addr()

			if err == nil {
				t.Fatal("succeed in p.Addr()", hostname, port)
			}

			if msg := err.Error(); msg != "invalid scheme: invalid" {
				t.Fatalf("unexpected error message: %#v", msg)
			}
		})
	})
}

func runHTTPServer(handler func(http.ResponseWriter, *http.Request)) (*httptest.Server, *url.URL) {
	s := httptest.NewServer(http.HandlerFunc(handler))

	u, err := url.Parse(s.URL)
	if err != nil {
		s.Close()
		panic(err)
	}

	return s, u
}

func TestHTTPPingerPing(t *testing.T) {
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
			t.Run(strconv.Itoa(success), func(t *testing.T) {
				s, u := runHTTPServer(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(success)
				})
				defer s.Close()

				p := &HTTPPinger{url: u}
				err := p.Ping(context.Background())
				if err != nil {
					t.Fatalf("failed in p.Ping(): %+#v", err)
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
			t.Run(strconv.Itoa(fail), func(t *testing.T) {
				s, u := runHTTPServer(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(fail)
				})
				defer s.Close()

				p := &HTTPPinger{url: u}
				err := p.Ping(context.Background())
				if err == nil {
					t.Fatal("succeed in p.Ping()")
				}

				if msg := err.Error(); msg != fmt.Sprintf("failed in HTTP request by %d %s", fail, http.StatusText(fail)) {
					t.Fatalf("unexpected error message: %#v", msg)
				}
			})
		}
	})

	t.Run("Context", func(t *testing.T) {
		s, u := runHTTPServer(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Second)
			w.WriteHeader(http.StatusOK)
		})
		defer s.Close()

		p := &HTTPPinger{url: u}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := p.Ping(ctx)
		<-ctx.Done()

		if err == nil {
			t.Fatal("succeed in p.Ping")
		}

		if msg := err.Error(); msg != "failed in HTTP request: Head "+s.URL+": context canceled" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		p := &HTTPPinger{url: &url.URL{Opaque: "::"}}
		err := p.Ping(context.Background())

		if err == nil {
			t.Fatal("succeed in p.Ping()")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in creating HTTP request: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
