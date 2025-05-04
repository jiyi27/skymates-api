package model

import (
	"github.com/google/uuid"
	"time"
)

// Post 帖子模型
type Post struct {
	ID        uuid.UUID `json:"id"`      // 更改为 uuid.UUID 类型
	UserID    uuid.UUID `json:"user_id"` // 更改为 uuid.UUID 类型
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` // 添加更新时间
}
