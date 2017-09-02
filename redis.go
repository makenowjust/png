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

	result, err := client.Ping().Result()
	if err != nil {
		return errors.Wrap(err, "failed to ping to redis")
	}

	if result != "PONG" {
		return errors.Errorf("invalid redis response: %#v", result)
	}

	return nil
}
