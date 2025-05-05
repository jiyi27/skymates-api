package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 对应 password 字段，隐藏字段
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url,omitempty"` // 可为 NULL
	Role      string    `json:"role"`                 // 对应 ENUM('user', 'admin')
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
