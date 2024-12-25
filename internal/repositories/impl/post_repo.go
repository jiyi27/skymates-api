package impl

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"skymates-api/internal/types"
)

type PostgresPostRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPostRepository(pool *pgxpool.Pool) *PostgresPostRepository {
	return &PostgresPostRepository{pool: pool}
}

func (s *PostgresPostRepository) Create(post *types.Post) error {
	// 实现创建用户逻辑
	return nil
}

func (s *PostgresPostRepository) GetByID(id string) (*types.Post, error) {
	// 实现获取用户信息逻辑
	return nil, nil
}
