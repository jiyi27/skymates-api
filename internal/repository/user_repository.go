package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	servererrors "skymates-api/errors"
	"skymates-api/internal/model"
	"time"
)

// QueryType 定义查询类型
// 可选值: 按用户名, 邮箱或ID查询用户
type QueryType int

const (
	QueryByUsername QueryType = iota // 按用户名查询
	QueryByEmail                     // 按邮箱查询
	QueryByID                        // 按ID查询
)

// UserRepository 定义用户存储库接口
type UserRepository interface {
	Create(user *model.User) error
	GetUserBy(queryType QueryType, value string) (*model.User, error)
	CheckExists(queryType QueryType, value string) (bool, error)
}

// MySQLUserRepository 实现了 UserRepository 接口, 使用 MySQL 数据库
type MySQLUserRepository struct {
	db *sqlx.DB // 使用 sqlx.DB 替代 sql.DB
}

// NewUserRepository 返回一个基于 MySQL 的用户存储库
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &MySQLUserRepository{db: db}
}

// Create 创建新用户
// 如果用户ID为空, 会生成一个新的 UUID
// 设置创建和更新时间, 并插入到 users 表
// 如果用户名或邮箱已存在, 返回 AlreadyExistsError
func (r *MySQLUserRepository) Create(user *model.User) error {
	// 处理时间戳
	now := time.Now()
	// 如果模型没有设置创建时间, 则填充
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	// 执行插入操作
	query := `INSERT INTO users (username, hashed_password, email, avatar_url, created_at, updated_at)
		VALUES (:username, :hashed_password, :email, :avatar_url, :created_at, :updated_at)`
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return servererrors.NewInternalError("创建用户失败", err)
	}
	return nil
}

// GetUserBy 根据查询类型和值检索用户
// 支持按用户名, 邮箱或ID查询
// 如果未找到, 返回 NotFoundError
// 查询失败时, 返回 InternalError
func (r *MySQLUserRepository) GetUserBy(queryType QueryType, value string) (*model.User, error) {
	var query string
	switch queryType {
	case QueryByUsername:
		query = `SELECT id, username, hashed_password, email, avatar_url, created_at, updated_at
			FROM users WHERE username = ?`
	case QueryByEmail:
		query = `SELECT id, username, hashed_password, email, avatar_url, created_at, updated_at
			FROM users WHERE email = ?`
	case QueryByID:
		query = `SELECT id, username, hashed_password, email, avatar_url, created_at, updated_at
			FROM users WHERE id = ?`
	default:
		return nil, servererrors.NewInternalError("无效的查询类型", nil)
	}

	var user model.User
	err := r.db.Get(&user, query, value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, servererrors.NewNotFoundError("用户未找到", err)
		}
		return nil, servererrors.NewInternalError("查询用户失败", err)
	}
	return &user, nil
}

// CheckExists 检查用户是否存在
// 使用 COUNT(1) 判断是否有匹配记录
// 如果查询失败, 返回 InternalError
func (r *MySQLUserRepository) CheckExists(queryType QueryType, value string) (bool, error) {
	var query string
	switch queryType {
	case QueryByUsername:
		query = `SELECT COUNT(1) FROM users WHERE username = ?`
	case QueryByEmail:
		query = `SELECT COUNT(1) FROM users WHERE email = ?`
	default:
		return false, servererrors.NewInternalError("无效的查询类型", nil)
	}

	var count int
	err := r.db.Get(&count, query, value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, servererrors.NewInternalError("检查用户存在性失败", err)
	}
	return count > 0, nil
}
