package model

import (
	"github.com/google/uuid"
	"time"
)

// Term 术语基础信息，用于列表展示多个术语
type Term struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// TermDetail 术语详细信息，用于详情展示
type TermDetail struct {
	Term
	// source, video_url, explanation 在数据库的定义都是可选(nullable)字段，所以使用指针类型
	// 使用指针类型，使其可以被被赋值为 nil，因为在数据库查询时，PostgreSQL 的 NULL 值会自动映射到 Go 的指针类型的 nil
	Explanation *string        `json:"explanation"`
	Source      *string        `json:"source"`
	VideoURL    *string        `json:"video_url"`
	Categories  []TermCategory `json:"categories"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TermCategory 术语分类
type TermCategory struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id"` // 当 ParentID 为 nil 时会序列化为 null
	CreatedAt time.Time  `json:"created_at"`
}
