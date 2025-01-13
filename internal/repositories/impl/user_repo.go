package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	servererrors "skymates-api/internal/errors"
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
	query := `INSERT INTO users (username, hashed_password, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := p.pool.QueryRow(context.Background(), query,
		user.Username,
		user.HashedPassword,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return servererrors.NewDatabaseError(fmt.Sprintf("PostgresUserRepository.Create: failed to insert user (username=%s, email=%s)",
			user.Username, user.Email), err)
	}

	return nil
}

// GetUserBy 通用查询函数
func (p *PostgresUserRepository) GetUserBy(field repositories.QueryField, value string) (*types.User, error) {
	// 构建查询语句
	query := `
        SELECT id, username, hashed_password, email
        FROM users
        WHERE ` + string(field) + ` = $1`

	user := &types.User{}
	err := p.pool.QueryRow(
		context.Background(),
		query,
		value,
	).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Email,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, servererrors.NewNotFoundError(
				fmt.Sprintf("PostgresUserRepository.GetUserBy: no such user, %s: %v", field, value),
				nil,
			)
		}
		return nil, servererrors.NewDatabaseError("PostgresUserRepository.GetUserBy: database error", err)
	}

	return user, nil
}

func (p *PostgresUserRepository) CheckExists(field repositories.QueryField, value string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM users WHERE %s = $1)`, field)

	var exists bool
	err := p.pool.QueryRow(context.Background(), query, value).Scan(&exists)

	if err != nil {
		return false, servererrors.NewDatabaseError(
			fmt.Sprintf("PostgresUserRepository.CheckExists: database error when checking %s", field),
			err,
		)
	}

	return exists, nil
}
