package repositories

import (
	"context"
	"github.com/google/uuid"
	"skymates-api/internal/types"
)

type QueryField string

const (
	QueryByID       QueryField = "id"
	QueryByUsername QueryField = "username"
	QueryByEmail    QueryField = "email"
)

type UserRepository interface {
	Create(user *types.User) error
	GetUserBy(field QueryField, value string) (*types.User, error)
	CheckExists(field QueryField, value string) (bool, error)
}

type TermRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*types.TermDetail, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID, lastID *uuid.UUID, limit int) ([]types.Term, bool, error)
	GetCategoryTermCount(ctx context.Context, categoryID uuid.UUID) (int, error)
}

type PostRepository interface {
	Create(post *types.Post) error
	GetByID(id string) (*types.Post, error)
}

type CommentRepository interface {
	Create(comment *types.Comment) error
}

//type Database interface {
//	CreateUser(user *types.User) error
//	GetUserByID(id string) (*types.User, error)
//	GetUserByUsername(username string) (*types.User, error)
//	UpdateUser(user *types.User) error
//
//	CreatePost(post *types.Post) error
//	GetPostByID(id string) (*types.Post, error)
//	ListPosts(offset, limit int) ([]*types.Post, error)
//	UpdatePost(post *types.Post) error
//	DeletePost(id string) error
//
//	CreateComment(comment *types.Comment) error
//	GetCommentsByPostID(postID string) ([]*types.Comment, error)
//	DeleteComment(id string) error
//
//	Close()
//	Ping() error
//}
