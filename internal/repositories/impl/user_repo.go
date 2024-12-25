package impl

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"skymates-api/internal/repositories"
	"skymates-api/internal/types"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) repositories.UserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (s *PostgresUserRepository) Create(user *types.User) error {
	// 实现创建用户逻辑
	return nil
}

func (s *PostgresUserRepository) GetByID(id string) (*types.User, error) {
	// 实现获取用户信息逻辑
	return nil, nil
}

func (s *PostgresUserRepository) GetByUsername(username string) (*types.User, error) {
	// 实现获取用户信息逻辑
	return nil, nil
}
