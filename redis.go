package png

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type RedisPinger struct {
	addr     string
	password string
	db       int
}

func (p *RedisPinger) Ping(ctx context.Context) error {
	opts := &redis.Options{
		Addr:     p.addr,
		Password: p.password,
		DB:       p.db,
	}

	if deadline, ok := ctx.Deadline(); ok {
		d := time.Until(deadline)
		opts.DialTimeout = d
		opts.ReadTimeout = d
		opts.WriteTimeout = d
	}

	client := redis.NewClient(opts)
	defer client.Close()

	done := make(chan error)
	go func() {
		result, err := client.Ping().Result()
		if err != nil {
			done <- errors.Wrap(err, "failed in PING command")
			return
		}

		if result != "PONG" {
			done <- errors.Errorf("invalid redis response: %#v", result)
			return
		}

		done <- nil
	}()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "failed in PING command")
	case err := <-done:
		return err
	}
}
