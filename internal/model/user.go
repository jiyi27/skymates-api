package model

import (
	"github.com/google/uuid"
	"time"
)

// User 用户模型
type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`          // json:"-" 确保不会在 JSON 响应中返回
	Email          string    `json:"email"`      // 修正Email字段名称的大小写
	AvatarURL      string    `json:"avatar_url"` // 修正字段名大小写
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
