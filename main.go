package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go-redis-essentials/database"
	"go-redis-essentials/user"
)

func main() {
	// 初始化 Redis 客戶端
	client, err := database.NewClient()
	if err != nil {
		fmt.Println("Redis 連線失敗:", err)
		return
	}

	// 建立使用者資料
	u := user.User{Name: "Elliot", Email: "elliot@mail.com", Age: 25}

	// 序列化：將 struct 轉為 JSON
	data, err := json.Marshal(u)
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
	var userFromRedis user.User
	err = json.Unmarshal([]byte(val), &userFromRedis)
	if err != nil {
		fmt.Println("反序列化錯誤:", err)
		return
	}
	fmt.Printf("反序列化後, 讀取到的使用者: %+v\n", userFromRedis)
}
