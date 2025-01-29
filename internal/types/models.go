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
	HashedPassword string    `json:"-"` // json:"-" 确保不会在 JSON 响应中返回
	Email          string    `json:"Email"`
	// omitempty 如果字段的值是其类型的零值, 则在 JSON 输出中省略该字段
	// omitempty 只对序列化有效, 不会对反序列化产生影响
	// 为了方便客户端处理, 这里我们不使用 omitempty
	// AvatarUrl string    `json:"avatar_url,omitempty"`
	AvatarUrl string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Term 术语基础信息, 用于列表展示多个术语
type Term struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// TermDetail 术语详细信息, 用于详情展示
type TermDetail struct {
	Term
	// source, video_url, explanation 在数据库的定义都是可选(nullable)字段, 所以使用指针类型
	// 使用指针类型, 使其可以被被赋值为 nil, 因为在数据库查询时, PostgreSQL 的 NULL 值会自动映射到 Go 的指针类型的 nil
	// 如果不使用指针类型, 则会报错, 无法将 NULL 值映射到非指针类型
	// *string 类型序列化为 JSON 时, 会序列化指针指向的值, 而不是地址
	Explanation *string        `json:"explanation"`
	Source      *string        `json:"source"`
	VideoURL    *string        `json:"video_url"`
	Categories  []TermCategory `json:"categories"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
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
