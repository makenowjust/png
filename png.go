package png

import (
	"context"
)

type Pinger interface {
	Addr() (hostname string, port int, err error)
	Ping(ctx context.Context) error
}
