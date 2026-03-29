package database

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

// NewClient 初始化並返回一個 Redis 客戶端實例
func NewClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 伺服器位址
		Password: "",               // 無密碼設定
		DB:       0,                // 使用預設資料庫
	})

	// 測試 Redis 連線
	pong, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("無法連線到 Redis: %w", err)
	}
	fmt.Println(pong) // 成功連線時印出 PONG
	return client, nil
}
