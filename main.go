package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("Go Redis Tutorial")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping(client.Context()).Result()
	fmt.Println(pong, err)
}
