package png

import (
	"context"
	"time"
)

type Timeout struct {
	Err error
}

func (timeout *Timeout) Error() string {
	return timeout.Err.Error()
}

type Pinger interface {
	Ping(ctx context.Context) error
}

func PingWithTimeout(p Pinger, timeout time.Duration) (elapsed time.Duration, err error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error)
	go func() {
		done <- p.Ping(ctx)
	}()

	select {
	case <-ctx.Done():
		elapsed = time.Since(start)
		err = &Timeout{Err: ctx.Err()}
	case err = <-done:
		elapsed = time.Since(start)
	}

	return
}
