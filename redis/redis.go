package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)
var client *redis.Client
func ConnectToRedis() {
	fmt.Println("Connecting to redis")
	client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	fmt.Println("Connected to redis client")
}

func GetRedisInstance() *redis.Client{
	return client
}
