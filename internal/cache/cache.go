package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	// Test Connection 
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Resis: " +err.Error())
	}
}
