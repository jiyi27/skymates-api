package impl

import (
	"skymates-api/internal/service"
	"skymates-api/internal/types"
)

type CommentService struct {
	db service.Database
}

func NewCommentService(db service.Database) *CommentService {
	return &CommentService{
		db: db,
	}
}

func (s *CommentService) CreateComment(postID, userID, content string) error {
	// 实现创建评论逻辑
	return nil
}

func (s *CommentService) ListComments(postID string) ([]*types.Comment, error) {
	// 实现评论列表逻辑
	return nil, nil
}
