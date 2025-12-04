package service

import (
	"context"
	"log"
	"micro-site/config"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RedisClient *redis.Client

func InitRedis() {
	cfg := config.LoadConfig()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost, // example: "localhost:6379"
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected!")
}
