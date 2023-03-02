package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fnmzgdt/e_shop/src/utils"
	"github.com/go-redis/redis/v8"
)

type CacheConnection struct {
	client *redis.Client
}

func SetupCacheConnection() (*CacheConnection, error) {
	var (
		redis_db, err = strconv.Atoi(utils.GetEnv("REDIS_DB_ID", ""))
		password      = utils.GetEnv("REDIS_PASSWORD", "")
		host          = utils.GetEnv("REDIS_HOST", "localhost:6379")
	)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       redis_db,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	fmt.Println("Successful conneciton to Redis.")
	return &CacheConnection{client}, nil
}

func (r *CacheConnection) GetKey(key string) (string, error) {
	ctx := context.Background()
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *CacheConnection) SetKey(key string, value interface{}, exp time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, exp).Err()
}
