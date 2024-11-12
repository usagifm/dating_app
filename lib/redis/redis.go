package redis

import (
	"context"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/usagifm/dating-app/lib/logger"
)

func Init(ctx context.Context, addr, password string) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	}

	redisClient := redis.NewClient(opts)
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		logger.GetLogger(ctx).Errorf("init redis fail: ", err)
		return nil, err
	}
	//redisClient.AddHook(nrredis.NewHook(opts))
	return redisClient, nil
}

func InitMockRedis() (*redis.Client, error) {
	mr, err := miniredis.Run()
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return client, err
}
