package impl

import (
	"skymates-api/internal/service"
	"skymates-api/internal/types"
)

type PostService struct {
	db service.Database
}

func NewPostService(db service.Database) *PostService {
	return &PostService{
		db: db,
	}
}

func (s *PostService) CreatePost(userID, content string) error {
	// 实现创建帖子逻辑
	return nil
}

func (s *PostService) GetPost(id string) (*types.Post, error) {
	// 实现获取帖子逻辑
	return nil, nil
}

func (s *PostService) ListPosts(page, size int) ([]*types.Post, error) {
	// 实现帖子列表逻辑
	return nil, nil
}
