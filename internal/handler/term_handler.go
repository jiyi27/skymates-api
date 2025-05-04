package handler

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	servererrors "skymates-api/internal/errors"
	"skymates-api/internal/middleware"
	"skymates-api/internal/repositories"
	"skymates-api/internal/types"
	"strconv"
)

type TermHandler struct {
	BaseHandler
	termRepo repositories.TermRepository
}

func RegisterTermRoutes(tr repositories.TermRepository, mux *http.ServeMux) {
	h := &TermHandler{termRepo: tr}
	// TODO: 为什么没有 匹配 OPTIONS 方法的路由, 却可以在前端发送 OPTIONS 请求并得到响应??
	// 如果 /api/terms/{id} 会覆盖 /api/terms/suggestions, 那按理说依然不应该出现 OPTIONS 请求得不到响应的情况, 因为都用到了 middleware.CORS(nil)
	mux.HandleFunc("/api/term/{id}", middleware.Use(h.GetTermDetail, middleware.Logger, middleware.CORS(nil)))
	mux.HandleFunc("/api/terms/suggestions", middleware.Use(h.GetTermSuggestions, middleware.Logger, middleware.CORS(nil)))
	mux.HandleFunc("/api/terms", middleware.Use(h.ListTermsByCategory, middleware.Logger, middleware.CORS(nil)))
}

func (h *TermHandler) GetTermSuggestions(w http.ResponseWriter, r *http.Request) {
	// 当 query string 中没有 query 时, 会返回空字符串
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	suggestions, err := h.termRepo.GetSuggestions(r.Context(), query)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.SendJSON(w, http.StatusOK, "Success", suggestions)
}

func (h *TermHandler) GetTermDetail(w http.ResponseWriter, r *http.Request) {
	termID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, "Invalid term ID", nil)
		return
	}

	term, err := h.termRepo.GetByID(r.Context(), termID)
	if err != nil {
		var serverErr *servererrors.ServerError
		if errors.As(err, &serverErr) {
			switch serverErr.Kind {
			case servererrors.KindNotFound:
				h.SendJSON(w, http.StatusNotFound, "Term not found", nil)
			default:
				h.SendJSON(w, http.StatusInternalServerError, "Internal server error", nil)
				log.Printf("TermHandler.GetTermDetail: failed to get term: %v", err)
			}
			return
		}

		h.SendJSON(w, http.StatusInternalServerError, "Internal server error", nil)
		log.Printf("TermHandler.GetTermDetail: failed to get term: %v", err)
		return
	}

	h.SendJSON(w, http.StatusOK, "Success", term)
}

// ListTermsByCategory 获取分类下的术语列表
func (h *TermHandler) ListTermsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryId, err := uuid.Parse(r.URL.Query().Get("categoryId"))
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	// 当 query string 中没有 last_id 时, 会返回空字符串
	// 当 last_id 为空时, lastID 为 nil
	lastIDStr := r.URL.Query().Get("last_id")
	var lastID *uuid.UUID
	if lastIDStr != "" {
		parsedID, err := uuid.Parse(lastIDStr)
		if err != nil {
			h.SendJSON(w, http.StatusBadRequest, "Invalid last_id format", nil)
			return
		}
		lastID = &parsedID
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20 // 默认值
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			h.SendJSON(w, http.StatusBadRequest, "Invalid limit value", nil)
			return
		}
		if parsedLimit > 50 {
			limit = 50
		} else {
			limit = parsedLimit
		}
	}

	terms, hasMore, err := h.termRepo.ListByCategory(r.Context(), categoryId, lastID, limit)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, "Internal server error", nil)
		log.Printf("TermHandler.ListTermsByCategory: failed to list terms: %v", err)
		return
	}

	var newLastID *uuid.UUID
	if len(terms) > 0 {
		newLastID = &terms[len(terms)-1].ID
	}

	response := types.ListTermsResponse{
		Terms:   terms,
		LastID:  newLastID,
		HasMore: hasMore,
	}

	h.SendJSON(w, http.StatusOK, "Success", response)
}
