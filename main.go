package main

import (
	"context"
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

	// 設定數值 (Set), 參數分別為：上下文 (ctx), 鍵 (key), 值 (value), 過期時間 (expiration)
	err = client.Set(context.Background(), "name", "WFFANG", 0).Err()
	if err != nil {
		fmt.Println("Error setting value:", err)
	}

	// 讀取數值 (Get)
	val, err := client.Get(context.Background(), "name").Result()
	if err != nil {
		fmt.Println("Error getting value:", err)
	} else {
		fmt.Println("Value:", val)
		fmt.Println("Key:", "name")
	}
}
