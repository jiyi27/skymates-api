package impl

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"skymates-api/internal/types"
)

type PostgresCommentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresCommentRepository(pool *pgxpool.Pool) *PostgresCommentRepository {
	return &PostgresCommentRepository{pool: pool}
}

func (s *PostgresCommentRepository) Create(comment *types.Comment) error {
	// 实现创建用户逻辑
	return nil
}
