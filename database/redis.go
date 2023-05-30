package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis(host, password string, port, db int) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis启动失败")
	}
}
