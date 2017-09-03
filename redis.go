package png

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type RedisPinger struct {
	hostname string
	port     int
	password string
	db       int
}

func (p *RedisPinger) Addr() (string, int, error) {
	return p.hostname, p.port, nil
}

func (p *RedisPinger) Ping(ctx context.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", p.hostname, p.port),
		Password: p.password,
		DB:       p.db,
	})
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
