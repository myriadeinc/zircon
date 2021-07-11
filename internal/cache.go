package internal

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

// var ctx = context.Background()
var Client *redis.Client

func CacheInit(ctx context.Context) {
	Client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Set(ctx context.Context, key string, value string) error {
	return Client.Set(ctx, key, value, 0).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}
