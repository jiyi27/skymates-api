package types

import (
	"github.com/google/uuid"
	"time"
)

/*
User
Tag 是元数据，用于指导序列化和反序列化过程
序列化是将 Go 语言中的数据结构（比如 struct、map、slice 等）
转换成一种可以存储或传输的格式，通常是 JSON 字符串或字节流
反序列化则是将这种存储或传输的格式（比如 JSON 字符串）转换回 Go 语言中的数据结构
json:"id" 显然是用于指导 JSON 序列化和反序列化的
*/
type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"` // 使用 json:"-" 确保不会在 JSON 响应中返回
	Email          string    `json:"email"`
	AvatarUrl      string    `json:"avatar,omitempty"` // 如果字段的值是其类型的零值, 则在 JSON 输出中省略该字段
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Term 术语基础信息, 用于列表展示多个术语
type Term struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// TermDetail 术语详细信息, 用于详情展示
type TermDetail struct {
	Term
	TextExplanation string         `json:"text_explanation"`
	Source          string         `json:"source"`
	VideoURL        string         `json:"video_url,omitempty"`
	Categories      []TermCategory `json:"categories"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type TermCategory struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id"` // 当 ParentID 为 nil 时会序列化为 null
	CreatedAt time.Time  `json:"created_at"`
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
