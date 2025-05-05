package service

import (
	"context"
	"log"
	servererrors "skymates-api/errors"
	"skymates-api/internal/model"
	"skymates-api/internal/repository"
)

// TermService 定义术语相关的业务逻辑接口
type TermService interface {
	SearchTerms(ctx context.Context, keyword string) ([]model.TermSummary, error)
	GetTermByID(ctx context.Context, id int64) (*model.TermDetail, error)
	ListTermsByCategory(ctx context.Context, categoryID int64, lastID *int64, limit int) ([]model.TermSummary, bool, error)
	CreateTerm(ctx context.Context, term *model.Term, categoryIDs []int64) (int64, error)
	UpdateTerm(ctx context.Context, term *model.Term, categoryIDs []int64) error
}

// termService 实现 TermService 接口
type termService struct {
	termRepository repository.TermRepository
}

// NewTermService 创建 TermService 实例
func NewTermService(termRepository repository.TermRepository) TermService {
	return &termService{
		termRepository: termRepository,
	}
}

// SearchTerms 根据关键字搜索术语
func (s *termService) SearchTerms(ctx context.Context, keyword string) ([]model.TermSummary, error) {
	terms, err := s.termRepository.SearchTerms(ctx, keyword)
	if err != nil {
		log.Printf("TermService.SearchTerms: %v", err)
		return nil, servererrors.NewInternalError("搜索术语失败", err)
	}
	summaries := make([]model.TermSummary, len(terms))
	for i, term := range terms {
		summaries[i] = model.TermSummary{ID: term.ID, Name: term.Name}
	}
	return summaries, nil
}

// GetTermByID 根据 ID 获取术语详情
func (s *termService) GetTermByID(ctx context.Context, id int64) (*model.TermDetail, error) {
	term, err := s.termRepository.GetTermByID(ctx, id)
	if err != nil {
		log.Printf("TermService.GetTermByID: %v", err)
		return nil, servererrors.NewInternalError("获取术语详情失败", err)
	}
	return term, nil
}

// ListTermsByCategory 列出指定分类下的术语
func (s *termService) ListTermsByCategory(ctx context.Context, categoryID int64, lastID *int64, limit int) ([]model.TermSummary, bool, error) {
	terms, hasMore, err := s.termRepository.ListTermsByCategory(ctx, categoryID, lastID, limit)
	if err != nil {
		log.Printf("TermService.ListTermsByCategory: %v", err)
		return nil, false, servererrors.NewInternalError("列出分类下的术语失败", err)
	}
	summaries := make([]model.TermSummary, len(terms))
	for i, term := range terms {
		summaries[i] = model.TermSummary{ID: term.ID, Name: term.Name}
	}
	return summaries, hasMore, nil
}

// CreateTerm 创建术语并关联分类
func (s *termService) CreateTerm(ctx context.Context, term *model.Term, categoryIDs []int64) (int64, error) {
	id, err := s.termRepository.CreateTerm(ctx, term, categoryIDs)
	if err != nil {
		log.Printf("TermService.CreateTerm: %v", err)
		return 0, servererrors.NewInternalError("创建术语失败", err)
	}
	return id, nil
}

// UpdateTerm 更新术语并更新关联分类
func (s *termService) UpdateTerm(ctx context.Context, term *model.Term, categoryIDs []int64) error {
	err := s.termRepository.UpdateTerm(ctx, term, categoryIDs)
	if err != nil {
		log.Printf("TermService.UpdateTerm: %v", err)
		return servererrors.NewInternalError("更新术语失败", err)
	}
	return nil
}
