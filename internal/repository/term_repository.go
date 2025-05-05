package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	servererrors "skymates-api/errors"
	"skymates-api/internal/types"
)

// TermRepository 定义术语存储库接口
type TermRepository interface {
	// SearchTerms 根据关键字搜索术语，最多返回 20 条记录
	// 优先返回完全匹配，其次前缀匹配，最后包含匹配；同一匹配级别下，名称更短者靠前
	SearchTerms(ctx context.Context, keyword string) ([]types.Term, error)
	// GetTermByID 根据 ID 查询术语详情，如果不存在则返回 NotFoundError
	GetTermByID(ctx context.Context, id int64) (*types.TermDetail, error)
	// ListTermsByCategory 列出指定分类下的术语，支持基于游标的分页
	// lastID 为上一次查询最后一条记录的 ID；limit 为每页大小
	// 返回值 hasMore 表示是否还有更多记录
	ListTermsByCategory(ctx context.Context, categoryID int64, lastID *int64, limit int) ([]types.Term, bool, error)
	// CountTermsInCategory 统计指定分类下的术语总数
	CountTermsInCategory(ctx context.Context, categoryID int64) (int, error)
}

// MySQLTermRepository 实例包含一个连接池 *sql.DB
type MySQLTermRepository struct {
	db *sql.DB
}

// NewTermRepository 创建并返回一个 MySQLTermRepository 实例
func NewTermRepository(db *sql.DB) TermRepository {
	return &MySQLTermRepository{db: db}
}

// SearchTerms 根据关键字搜索术语，最多返回 20 条记录
// 优先返回完全匹配，其次前缀匹配，最后包含匹配；同一匹配级别下，名称更短者靠前
func (r *MySQLTermRepository) SearchTerms(ctx context.Context, keyword string) ([]types.Term, error) {
	// 构建模糊匹配模式
	exactPattern := keyword                // 完全匹配
	prefixPattern := keyword + "%"         // 前缀匹配
	containsPattern := "%" + keyword + "%" // 包含匹配

	query := `
		SELECT id, name
		FROM terms
		WHERE name LIKE ?
		ORDER BY
		  CASE
		    WHEN name = ? THEN 1
		    WHEN name LIKE ? THEN 2
		    ELSE 3
		  END,
		  CHAR_LENGTH(name)
		LIMIT 20`

	// 执行查询
	rows, err := r.db.QueryContext(ctx, query,
		containsPattern,
		exactPattern,
		prefixPattern,
	)
	if err != nil {
		return nil, servererrors.NewInternalError("搜索术语失败", err)
	}
	defer rows.Close()

	// 解析结果集
	var terms []types.Term
	for rows.Next() {
		var t types.Term
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, servererrors.NewInternalError("扫描术语失败", err)
		}
		terms = append(terms, t)
	}
	if err := rows.Err(); err != nil {
		return nil, servererrors.NewInternalError("迭代术语失败", err)
	}
	return terms, nil
}

// GetTermByID 根据 ID 查询术语详情，如果不存在则返回 NotFoundError
func (r *MySQLTermRepository) GetTermByID(ctx context.Context, id int64) (*types.TermDetail, error) {
	// 使用 JSON_ARRAYAGG 和 JSON_OBJECT 构建分类数组
	query := `
		SELECT
		  t.id,
		  t.name,
		  t.explanation,
		  t.source,
		  t.video_url,
		  t.created_at,
		  t.updated_at,
		  COALESCE(
		    JSON_ARRAYAGG(
		      JSON_OBJECT(
		        'id', c.id,
		        'name', c.name,
		        'parent_id', c.parent_id,
		        'created_at', c.created_at
		      )
		    ),
		    JSON_ARRAY()
		  ) AS categories
		FROM terms t
		LEFT JOIN term_category_relations tcr ON t.id = tcr.term_id
		LEFT JOIN term_categories c ON tcr.category_id = c.id
		WHERE t.id = ?
		GROUP BY t.id, t.name, t.explanation, t.source, t.video_url, t.created_at, t.updated_at
	`

	row := r.db.QueryRowContext(ctx, query, id)
	var detail types.TermDetail
	if err := row.Scan(
		&detail.ID,
		&detail.Name,
		&detail.Explanation,
		&detail.Source,
		&detail.VideoURL,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.Categories,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, servererrors.NewNotFoundError(fmt.Sprintf("未找到 ID=%d 的术语", id), err)
		}
		return nil, servererrors.NewInternalError("查询术语详情失败", err)
	}
	return &detail, nil
}

// ListTermsByCategory 列出指定分类下的术语，支持基于游标的分页
// lastID 为上一次查询最后一条记录的 ID；limit 为每页大小
// 返回值 hasMore 表示是否还有更多记录
func (r *MySQLTermRepository) ListTermsByCategory(ctx context.Context, categoryID int64, lastID *int64, limit int) ([]types.Term, bool, error) {
	// 多取一条用于判断是否有更多
	query := `
		SELECT t.id, t.name
		FROM terms t
		JOIN term_category_relations tcr ON t.id = tcr.term_id
		WHERE tcr.category_id = ?
		  AND (? IS NULL OR t.id < ?)
		ORDER BY t.id DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query,
		categoryID,
		lastID,
		lastID,
		limit+1,
	)
	if err != nil {
		return nil, false, servererrors.NewInternalError("列出分类术语失败", err)
	}
	defer rows.Close()

	var terms []types.Term
	for rows.Next() {
		var t types.Term
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, false, servererrors.NewInternalError("扫描术语失败", err)
		}
		terms = append(terms, t)
	}
	if err := rows.Err(); err != nil {
		return nil, false, servererrors.NewInternalError("迭代分类术语失败", err)
	}

	// 判断是否有更多
	hasMore := len(terms) > limit
	if hasMore {
		terms = terms[:limit]
	}
	return terms, hasMore, nil
}

// CountTermsInCategory 统计指定分类下的术语总数
func (r *MySQLTermRepository) CountTermsInCategory(ctx context.Context, categoryID int64) (int, error) {
	query := `SELECT COUNT(*) FROM term_category_relations WHERE category_id = ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count); err != nil {
		return 0, servererrors.NewInternalError("统计分类术语数量失败", err)
	}
	return count, nil
}
