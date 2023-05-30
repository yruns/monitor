package tests

//import (
//	"context"
//	"fmt"
//	"github.com/redis/go-redis/v9"
//	"testing"
//	"time"
//)
//
//func TestRedis(t *testing.T) {
//	Redis := redis.NewClient(&redis.Options{
//		Addr:     fmt.Sprintf("%s:%d", "localhost", 6379),
//		Password: "",
//		DB:       0,
//	})
//
//	_, err := Redis.Ping(context.Background()).Result()
//	if err != nil {
//		fmt.Println("Redis启动失败")
//	}
//
//	Redis.SetEx(context.Background(), "1234", "sjlg", time.Minute*3)
//}
