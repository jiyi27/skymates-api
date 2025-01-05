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
		return servererrors.NewDatabaseError(fmt.Sprintf("repository.User.Create: failed to insert user (username=%s, email=%s)",
			user.Username, user.Email), err)
	}

	return nil
}

func (p *PostgresUserRepository) GetByID(id string) (*types.User, error) {
	query := `
	   SELECT id, username, hashed_password, email, created_at, updated_at
	   FROM users
	   WHERE id = $1`

	user := &types.User{}
	err := p.pool.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, servererrors.NewNotFoundError(fmt.Sprintf("repository.User.GetByID: no such user, ID: %v", id), nil)
		}
		return nil, servererrors.NewDatabaseError("repository.User.GetByID: database error", err)
	}

	return user, nil
}

func (p *PostgresUserRepository) GetByUsername(username string) (*types.User, error) {
	query := `
	   SELECT id, username, hashed_password, email, created_at, updated_at
	   FROM users
	   WHERE username = $1`

	user := &types.User{}
	err := p.pool.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, servererrors.NewNotFoundError(fmt.Sprintf("repository.User.GetByUsername: no such user, username: %v", username), nil)
		}
		return nil, servererrors.NewDatabaseError("repository.User.GetByUsername: database error", err)
	}

	return user, nil
}

func (p *PostgresUserRepository) CheckUsernameExists(username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := p.pool.QueryRow(context.Background(), query, username).Scan(&exists)

	if err != nil {
		return false, servererrors.NewDatabaseError("repository.User.CheckUsernameExists: database error", err)
	}

	return exists, nil
}

func (p *PostgresUserRepository) CheckEmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := p.pool.QueryRow(context.Background(), query, email).Scan(&exists)

	if err != nil {
		return false, servererrors.NewDatabaseError("repository.User.CheckEmailExists: database error", err)
	}

	return exists, nil
}
