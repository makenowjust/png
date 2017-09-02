package png

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type RedisPinger struct {
	Addr     string
	Password string
	DB       int
}

func (pinger *RedisPinger) Ping(ctx context.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     pinger.Addr,
		Password: pinger.Password,
		DB:       pinger.DB,
	})
	defer client.Close()

	done := make(chan error)
	go func() {
		result, err := client.Ping().Result()
		if err != nil {
			done <- errors.Wrap(err, "failed to ping to redis")
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
		return errors.Wrap(ctx.Err(), "failed to ping to redis")
	case err := <-done:
		return err
	}
}
