package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	// 初始化 Redis 客戶端
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 伺服器位址
		Password: "",               // 無密碼設定
		DB:       0,                // 使用預設資料庫
	})

	// 測試 Redis 連線
	pong, err := client.Ping(client.Context()).Result()
	fmt.Println(pong, err)

	// 建立使用者資料
	user := User{Name: "Elliot", Email: "elliot@mail.com", Age: 25}

	// 序列化：將 struct 轉為 JSON
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("序列化錯誤:", err)
		return
	}

	// 建立 context
	ctx := context.Background()

	// 存入 Redis（過期時間設為 0 表示永不過期）
	err = client.Set(ctx, "user:1", data, 0).Err()
	if err != nil {
		fmt.Println("存入 Redis 錯誤:", err)
		return
	}

	// 從 Redis 取得資料
	val, err := client.Get(ctx, "user:1").Result()
	if err != nil {
		fmt.Println("從 Redis 讀取錯誤:", err)
		return
	}
	fmt.Println("原始 JSON 字串:", val)

	// 反序列化：將 JSON 轉回 struct
	var userFromRedis User
	err = json.Unmarshal([]byte(val), &userFromRedis)
	if err != nil {
		fmt.Println("反序列化錯誤:", err)
		return
	}
	fmt.Printf("反序列化後, 讀取到的使用者: %+v\n", userFromRedis)
}
