package v1

import "github.com/google/uuid"

// CreatePostRequest 创建帖子请求
type CreatePostRequest struct {
	Content string `json:"content"`
}

// PostResponse 帖子响应
type PostResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"` // 增加用户名称，方便前端展示
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"` // 格式化的时间字符串
}
