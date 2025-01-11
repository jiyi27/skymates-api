package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	servererrors "skymates-api/internal/errors"
	"skymates-api/internal/repositories"
	"skymates-api/internal/types"
)

type PostgresTermRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTermRepository(pool *pgxpool.Pool) repositories.TermRepository {
	return &PostgresTermRepository{pool: pool}
}

func (p *PostgresTermRepository) GetByID(ctx context.Context, id uuid.UUID) (*types.TermDetail, error) {
	query := `
		SELECT t.id, t.name, t.text_explanation, t.source, t.video_url, t.created_at, t.updated_at,
		       -- COALESCE 确保在没有分类时返回空数组而不是NULL
			   COALESCE(
			       -- postgresql json_agg() 聚合函数用于将行转换为 JSON 数组, 必须在 GROUP BY 子句中使用
				   json_agg(json_build_object('id', c.id,'name', c.name,'parent_id', c.parent_id, 'created_at', c.created_at)) 
				   -- FILTER 子句确保只聚合非 NULL 的分类记录
				   FILTER (WHERE c.id IS NOT NULL), 
				   '[]'
			   ) as categories
		FROM terms t
		LEFT JOIN term_category_relations tcr ON t.id = tcr.term_id
		LEFT JOIN term_categories c ON tcr.category_id = c.id
		WHERE t.id = $1
		-- 在 GROUP BY 之前, 得到的是每个分类的单独行, 即一个分类对应一行
		-- GROUP BY 之后, 将 t.id, t.name, t.text_explanation 等值相同的行合并为一行, 
		-- 此时 categories 自然就是一个数组了, categories 是上面 SELECT 中的别名
		GROUP BY t.id, t.name, t.text_explanation, t.source, t.video_url,
				 t.created_at, t.updated_at`

	// 执行 Query 语句
	var term types.TermDetail
	err := p.pool.QueryRow(ctx, query, id).Scan(
		&term.ID,
		&term.Name,
		&term.TextExplanation,
		&term.Source,
		&term.VideoURL,
		&term.CreatedAt,
		&term.UpdatedAt,
		&term.Categories,
	)

	// pgx.ErrNoRows 只会由 QueryRow.Scan() 返回, Query() 不会返回 ErrNoRows,
	// 因为对于 Query() 来说, 查询结果为空是一个有效的状态
	// 会返回一个空的 Rows 对象, 调用 Next() 会返回 false
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, servererrors.NewNotFoundError(fmt.Sprintf("PostgresTermRepository.GetByID: no such term, ID: %v", id), nil)
		}
		return nil, servererrors.NewDatabaseError("PostgresTermRepository.GetByID: database error", err)
	}

	return &term, nil
}

func (p *PostgresTermRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID, lastID *uuid.UUID, limit int) ([]types.Term, bool, error) {
	// offset 分页在数据量大时性能会下降，因为数据库需要跳过 offset 行数据
	// 在高并发下可能出现数据重复或遗漏, 因为有可能在两次查询之间有新数据插入
	// 使用 lastID 分页可以避免这个问题, lastID 是上次查询结果中最后一条数据的 ID, 像个游标 cursor 一样
	query := `
		SELECT t.id, t.name
		FROM terms t
		JOIN term_category_relations tcr ON t.id = tcr.term_id
		WHERE tcr.category_id = $1
		AND ($2::uuid IS NULL OR t.id < $2)
		ORDER BY t.id DESC  -- 按 ID 降序排列, 保证 lastID 之前的数据在前面
		LIMIT $3`

	// 故意多查询一条数据, 用于判断是否还有更多数据
	// pgx.ErrNoRows 只会由 QueryRow.Scan() 返回, Query() 不会返回 ErrNoRows,
	// 因为对于 Query() 来说, 查询结果为空是一个有效的状态
	// 会返回一个空的 Rows 对象, 调用 rows.Next() 会返回 false
	rows, err := p.pool.Query(ctx, query, categoryID, lastID, limit+1)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	terms := make([]types.Term, 0, limit+1)
	for rows.Next() {
		var term types.Term
		if err := rows.Scan(&term.ID, &term.Name); err != nil {
			return nil, false, err
		}
		terms = append(terms, term)
	}

	hasMore := len(terms) > limit
	if hasMore {
		terms = terms[:limit]
	}

	return terms, hasMore, nil
}

func (p *PostgresTermRepository) GetCategoryTermCount(ctx context.Context, categoryID uuid.UUID) (int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM term_category_relations
		WHERE category_id = $1
	`, categoryID).Scan(&count)

	return count, err
}