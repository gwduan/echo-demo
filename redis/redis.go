package redis

import (
	"context"
	"echo-demo/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func ClientInit() error {
	opt, err := redis.ParseURL(config.RedisURL())
	if err != nil {
		return err
	}

	cli := redis.NewClient(opt)

	ctx := context.Background()
	err = cli.Set(ctx, "tempFoo", "bar", time.Second).Err()
	if err != nil {
		cli.Close()
		return err
	}

	client = cli

	return nil
}

func Client() *redis.Client {
	return client
}
