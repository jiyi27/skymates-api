package impl

import (
	"context"
	"fmt"
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

func (p *PostgresUserRepository) Create(user *types.User) error {
	query := `INSERT INTO users (username, password, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := p.pool.QueryRow(context.Background(), query,
		user.Username,
		user.HashedPassword,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

func (p *PostgresUserRepository) GetByID(id string) (*types.User, error) {
	return nil, nil
	//// 根据用户ID查询用户信息
	//query := `
	//    SELECT id, username, password_hash, email, created_at, updated_at
	//    FROM users
	//    WHERE id = $1
	//`
	//
	//// 初始化用户结构体
	//user := &types.User{}
	//
	//// 执行查询并扫描结果到结构体
	//err := p.pool.QueryRow(
	//	context.Background(),
	//	query,
	//	id,
	//).Scan(
	//	&user.ID,
	//	&user.Username,
	//	&user.HashedPassword,
	//	&user.Email,
	//	&user.CreatedAt,
	//	&user.UpdatedAt,
	//)
	//
	//// 处理查询错误
	//if err != nil {
	//	if err == pgx.ErrNoRows {
	//		return nil, fmt.Errorf("未找到ID为%s的用户", id)
	//	}
	//	return nil, fmt.Errorf("查询用户失败: %w", err)
	//}
	//
	//return user, nil
}

func (p *PostgresUserRepository) GetByUsername(username string) (*types.User, error) {
	return nil, nil
	//// 根据用户名查询用户信息
	//query := `
	//    SELECT id, username, password_hash, email, created_at, updated_at
	//    FROM users
	//    WHERE username = $1
	//`
	//
	//// 初始化用户结构体
	//user := &types.User{}
	//
	//// 执行查询并扫描结果到结构体
	//err := p.pool.QueryRow(
	//	context.Background(),
	//	query,
	//	username,
	//).Scan(
	//	&user.ID,
	//	&user.Username,
	//	&user.HashedPassword,
	//	&user.Email,
	//	&user.CreatedAt,
	//	&user.UpdatedAt,
	//)
	//
	//// 处理查询错误
	//if err != nil {
	//	if err == pgx.ErrNoRows {
	//		return nil, fmt.Errorf("未找到用户名为%s的用户", username)
	//	}
	//	return nil, fmt.Errorf("查询用户失败: %w", err)
	//}
	//
	//return user, nil
}

func (p *PostgresUserRepository) CheckUsernameExists(username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := p.pool.QueryRow(context.Background(), query, username).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("检查用户名是否存在失败: %w", err)
	}

	return exists, nil
}
