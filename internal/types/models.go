package types

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"` // 使用 json:"-" 确保不会在 JSON 响应中返回
	Email          string    `json:"email"`
	AvatarUrl      string    `json:"avatar,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID      string `json:"id"`
	PostID  string `json:"post_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}
