package model

import (
	"github.com/google/uuid"
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uuid.UUID `json:"id"`      // 更改为 uuid.UUID 类型
	PostID    uuid.UUID `json:"post_id"` // 更改为 uuid.UUID 类型
	UserID    uuid.UUID `json:"user_id"` // 更改为 uuid.UUID 类型
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"` // 添加创建时间
	UpdatedAt time.Time `json:"updated_at"` // 添加更新时间
}
