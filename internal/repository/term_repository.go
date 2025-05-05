package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"skymates-api/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

// TermRepository 定义术语存储库接口
type TermRepository interface {
	SearchTerms(ctx context.Context, keyword string) ([]model.Term, error)
	GetTermByID(ctx context.Context, id int64) (*model.TermDetail, error)
	ListTermsByCategory(ctx context.Context, categoryID int64, lastID *int64, limit int) ([]model.Term, bool, error)
	CreateTerm(ctx context.Context, term *model.Term, categoryIDs []int64) (int64, error)
	UpdateTerm(ctx context.Context, term *model.Term, categoryIDs []int64) error
}

// TermRepositoryImpl 实现 TermRepository 接口
type TermRepositoryImpl struct {
	db *sqlx.DB
}

// NewTermRepository 创建 TermRepository 实例
func NewTermRepository(db *sqlx.DB) TermRepository {
	return &TermRepositoryImpl{db: db}
}

// SearchTerms 根据关键字搜索术语
func (r *TermRepositoryImpl) SearchTerms(ctx context.Context, keyword string) ([]model.Term, error) {
	query := `SELECT id, name FROM terms WHERE name LIKE ?`
	var terms []model.Term
	err := r.db.SelectContext(ctx, &terms, query, "%"+keyword+"%")
	if err != nil {
		log.Printf("TermRepositoryImpl.SearchTerms: %v", err)
		return nil, err
	}
	return terms, nil
}

// GetTermByID 根据 ID 获取术语详情
func (r *TermRepositoryImpl) GetTermByID(ctx context.Context, id int64) (*model.TermDetail, error) {
	query := `SELECT id, name, explanation, source_url, created_at, updated_at FROM terms WHERE id = ?`
	var term model.TermDetail
	err := r.db.GetContext(ctx, &term, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("TermRepositoryImpl.GetTermByID: %v", err)
		return nil, err
	}

	// 获取关联的分类 ID
	categoryQuery := `SELECT category_id FROM term_category_relations WHERE term_id = ?`
	var categoryIDs []int64
	err = r.db.SelectContext(ctx, &categoryIDs, categoryQuery, id)
	if err != nil {
		log.Printf("TermRepositoryImpl.GetTermByID: %v", err)
		return nil, err
	}
	term.CategoryIDs = categoryIDs
	return &term, nil
}

// ListTermsByCategory 列出指定分类下的术语
func (r *TermRepositoryImpl) ListTermsByCategory(ctx context.Context, categoryID int64, lastID *int64, limit int) ([]model.Term, bool, error) {
	var query string
	var args []interface{}
	if lastID == nil {
		query = `SELECT t.id, t.name FROM terms t
                 JOIN term_category_relations r ON t.id = r.term_id
                 WHERE r.category_id = ? ORDER BY t.id ASC LIMIT ?`
		args = []interface{}{categoryID, limit + 1}
	} else {
		query = `SELECT t.id, t.name FROM terms t
                 JOIN term_category_relations r ON t.id = r.term_id
                 WHERE r.category_id = ? AND t.id > ? ORDER BY t.id ASC LIMIT ?`
		args = []interface{}{categoryID, *lastID, limit + 1}
	}

	var terms []model.Term
	err := r.db.SelectContext(ctx, &terms, query, args...)
	if err != nil {
		log.Printf("TermRepositoryImpl.ListTermsByCategory: %v", err)
		return nil, false, err
	}

	hasMore := len(terms) > limit
	if hasMore {
		terms = terms[:limit]
	}
	return terms, hasMore, nil
}

// CreateTerm 创建新术语并关联分类
func (r *TermRepositoryImpl) CreateTerm(ctx context.Context, term *model.Term, categoryIDs []int64) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("TermRepositoryImpl.CreateTerm: %v", err)
		return 0, err
	}
	defer func(tx *sqlx.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Printf("TermRepositoryImpl.CreateTerm: %v", err)
		}
	}(tx)

	// 插入 terms 表
	query := `INSERT INTO terms (name, explanation, source_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, query, term.Name, term.Explanation, term.SourceURL, time.Now(), time.Now())
	if err != nil {
		log.Printf("TermRepositoryImpl.CreateTerm: %v", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("TermRepositoryImpl.CreateTerm: %v", err)
		return 0, err
	}

	// 插入 term_category_relations 表
	for _, categoryID := range categoryIDs {
		_, err := tx.ExecContext(ctx, `INSERT INTO term_category_relations (term_id, category_id) VALUES (?, ?)`, id, categoryID)
		if err != nil {
			log.Printf("TermRepositoryImpl.CreateTerm: %v", err)
			return 0, err
		}
	}

	// 更新 terms 表的 category_list 字段
	categoryList, _ := json.Marshal(categoryIDs)
	_, err = tx.ExecContext(ctx, `UPDATE terms SET category_list = ? WHERE id = ?`, categoryList, id)
	if err != nil {
		log.Printf("TermRepositoryImpl.CreateTerm: %v", err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("TermRepositoryImpl.CreateTerm: %v", err)
		return 0, err
	}
	return id, nil
}

// UpdateTerm 更新术语并更新关联分类
func (r *TermRepositoryImpl) UpdateTerm(ctx context.Context, term *model.Term, categoryIDs []int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("TermRepositoryImpl.UpdateTerm: %v", err)
		return err
	}
	defer func(tx *sqlx.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Printf("TermRepositoryImpl.UpdateTerm: %v", err)
		}
	}(tx)

	// 更新 terms 表
	query := `UPDATE terms SET name = ?, explanation = ?, source_url = ?, updated_at = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, query, term.Name, term.Explanation, term.SourceURL, time.Now(), term.ID)
	if err != nil {
		log.Printf("TermRepositoryImpl.UpdateTerm: %v", err)
		return err
	}

	// 删除旧的关联
	_, err = tx.ExecContext(ctx, `DELETE FROM term_category_relations WHERE term_id = ?`, term.ID)
	if err != nil {
		log.Printf("TermRepositoryImpl.UpdateTerm: %v", err)
		return err
	}

	// 插入新的关联
	for _, categoryID := range categoryIDs {
		_, err := tx.ExecContext(ctx, `INSERT INTO term_category_relations (term_id, category_id) VALUES (?, ?)`, term.ID, categoryID)
		if err != nil {
			log.Printf("TermRepositoryImpl.UpdateTerm: %v", err)
			return err
		}
	}

	// 更新 terms 表的 category_list 字段
	categoryList, _ := json.Marshal(categoryIDs)
	_, err = tx.ExecContext(ctx, `UPDATE terms SET category_list = ? WHERE id = ?`, categoryList, term.ID)
	if err != nil {
		log.Printf("TermRepositoryImpl.UpdateTerm: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("TermRepositoryImpl.UpdateTerm: %v", err)
		return err
	}
	return nil
}
