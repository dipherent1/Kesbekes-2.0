package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func CreateRedisClient() *redis.Client {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Adjust as needed
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	// Check if Redis connection works
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return client
}
