package service

import "skymates-api/internal/types"

type Database interface {
	// 定义数据库操作接口
	// 这里可以是具体的数据库实现，比如 MySQL、MongoDB 等
}

type UserService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	GetUser(id string) (*types.User, error)
}

type PostService interface {
	CreatePost(userID, content string) error
	GetPost(id string) (*types.Post, error)
	ListPosts(page, size int) ([]*types.Post, error)
}

type CommentService interface {
	CreateComment(postID, userID, content string) error
	ListComments(postID string) ([]*types.Comment, error)
}
