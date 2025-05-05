package handler

import (
	"errors"
	"log"
	"net/http"
	serverErrors "skymates-api/errors"
	v1 "skymates-api/internal/dto/v1"
	"skymates-api/internal/model"
	"skymates-api/internal/service"
	"skymates-api/internal/validator"
	"strconv"
)

// TermHandler 术语处理器
type TermHandler struct {
	BaseHandler
	termService service.TermService
}

// NewTermHandler 创建术语处理器
func NewTermHandler(termService service.TermService) *TermHandler {
	return &TermHandler{
		termService: termService,
	}
}

// SearchTerms 处理术语搜索请求
func (h *TermHandler) SearchTerms(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		h.ResponseJSON(w, http.StatusBadRequest, "缺少关键字", nil)
		return
	}

	terms, err := h.termService.SearchTerms(r.Context(), keyword)
	if err != nil {
		var serverErr *serverErrors.ServerError
		if errors.As(err, &serverErr) {
			h.ResponseJSON(w, http.StatusInternalServerError, serverErr.Message, nil)
			return
		}
		h.ResponseJSON(w, http.StatusInternalServerError, "服务器内部错误", nil)
		log.Printf("TermHandler.SearchTerms: %v", err)
		return
	}

	// 类型转换：model.TermSummary -> v1.TermSummary
	v1Terms := make([]v1.TermSummary, len(terms))
	for i, term := range terms {
		v1Terms[i] = v1.TermSummary{ID: term.ID, Name: term.Name}
	}

	response := v1.SearchTermsResponse{Terms: v1Terms}
	h.ResponseJSON(w, http.StatusOK, "成功", response)
}

// GetTermByID 处理获取术语详情请求
func (h *TermHandler) GetTermByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "无效的术语 ID", nil)
		return
	}

	term, err := h.termService.GetTermByID(r.Context(), id)
	if err != nil {
		var serverErr *serverErrors.ServerError
		if errors.As(err, &serverErr) {
			h.ResponseJSON(w, http.StatusInternalServerError, serverErr.Message, nil)
			return
		}
		h.ResponseJSON(w, http.StatusInternalServerError, "服务器内部错误", nil)
		log.Printf("TermHandler.GetTermByID: %v", err)
		return
	}

	response := v1.TermDetailResponse{
		ID:          term.ID,
		Name:        term.Name,
		Explanation: term.Explanation,
		SourceURL:   term.SourceURL,
		CategoryIDs: term.CategoryIDs,
		CreatedAt:   term.CreatedAt,
		UpdatedAt:   term.UpdatedAt,
	}
	h.ResponseJSON(w, http.StatusOK, "成功", response)
}

// ListTermsByCategory 处理列出分类下术语请求
func (h *TermHandler) ListTermsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.PathValue("categoryID")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "无效的分类 ID", nil)
		return
	}

	lastIDStr := r.URL.Query().Get("lastID")
	var lastID *int64
	if lastIDStr != "" {
		id, err := strconv.ParseInt(lastIDStr, 10, 64)
		if err != nil {
			h.ResponseJSON(w, http.StatusBadRequest, "无效的 lastID", nil)
			return
		}
		lastID = &id
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // 默认分页大小
	}

	terms, hasMore, err := h.termService.ListTermsByCategory(r.Context(), categoryID, lastID, limit)
	if err != nil {
		var serverErr *serverErrors.ServerError
		if errors.As(err, &serverErr) {
			h.ResponseJSON(w, http.StatusInternalServerError, serverErr.Message, nil)
			return
		}
		h.ResponseJSON(w, http.StatusInternalServerError, "服务器内部错误", nil)
		log.Printf("TermHandler.ListTermsByCategory: %v", err)
		return
	}

	// 类型转换：model.TermSummary -> v1.TermSummary
	v1Terms := make([]v1.TermSummary, len(terms))
	for i, term := range terms {
		v1Terms[i] = v1.TermSummary{ID: term.ID, Name: term.Name}
	}

	response := v1.ListTermsByCategoryResponse{
		Terms:   v1Terms,
		HasMore: hasMore,
	}
	h.ResponseJSON(w, http.StatusOK, "成功", response)
}

// CreateTerm 处理创建术语请求
func (h *TermHandler) CreateTerm(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateTermRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "请求格式无效", nil)
		return
	}

	msg, err := validator.ValidateRequest(req)
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, msg, nil)
		return
	}

	term := &model.Term{
		Name:        req.Name,
		Explanation: req.Explanation,
		SourceURL:   req.SourceURL,
	}

	id, err := h.termService.CreateTerm(r.Context(), term, req.CategoryIDs)
	if err != nil {
		var serverErr *serverErrors.ServerError
		if errors.As(err, &serverErr) {
			h.ResponseJSON(w, http.StatusInternalServerError, serverErr.Message, nil)
			return
		}
		h.ResponseJSON(w, http.StatusInternalServerError, "服务器内部错误", nil)
		log.Printf("TermHandler.CreateTerm: %v", err)
		return
	}

	h.ResponseJSON(w, http.StatusCreated, "术语创建成功", map[string]int64{"id": id})
}

// UpdateTerm 处理更新术语请求
func (h *TermHandler) UpdateTerm(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "无效的术语 ID", nil)
		return
	}

	var req v1.UpdateTermRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "请求格式无效", nil)
		return
	}

	msg, err := validator.ValidateRequest(req)
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, msg, nil)
		return
	}

	term := &model.Term{
		ID:          id,
		Name:        req.Name,
		Explanation: req.Explanation,
		SourceURL:   req.SourceURL,
	}

	err = h.termService.UpdateTerm(r.Context(), term, req.CategoryIDs)
	if err != nil {
		var serverErr *serverErrors.ServerError
		if errors.As(err, &serverErr) {
			h.ResponseJSON(w, http.StatusInternalServerError, serverErr.Message, nil)
			return
		}
		h.ResponseJSON(w, http.StatusInternalServerError, "服务器内部错误", nil)
		log.Printf("TermHandler.UpdateTerm: %v", err)
		return
	}

	h.ResponseJSON(w, http.StatusOK, "术语更新成功", nil)
}
