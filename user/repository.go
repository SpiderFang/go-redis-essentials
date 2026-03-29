package user

import (
	"context"
	"encoding/json"
	"fmt"

	// 確保這是您的 go module 名稱
	"github.com/go-redis/redis/v8"
)

// Repository 負責處理使用者資料的持久化
type Repository struct {
	client *redis.Client
}

// NewRepository 建立一個新的使用者倉庫
func NewRepository(client *redis.Client) *Repository {
	return &Repository{client: client}
}

// Save 將使用者資料存入 Redis
func (r *Repository) Save(ctx context.Context, id string, user *User) error {
	// 序列化：將 struct 轉為 JSON
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("序列化使用者錯誤: %w", err)
	}

	// 存入 Redis（過期時間設為 0 表示永不過期）
	key := fmt.Sprintf("user:%s", id)
	err = r.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("存入 Redis 錯誤: %w", err)
	}
	return nil
}

// GetByID 從 Redis 根據 ID 取得使用者資料
func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	key := fmt.Sprintf("user:%s", id)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("從 Redis 讀取錯誤: %w", err)
	}

	// 反序列化：將 JSON 轉回 struct
	var userFromRedis User
	err = json.Unmarshal([]byte(val), &userFromRedis)
	if err != nil {
		return nil, fmt.Errorf("反序列化使用者錯誤: %w", err)
	}

	return &userFromRedis, nil
}
