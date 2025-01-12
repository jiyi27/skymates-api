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
)

type TermHandler struct {
	BaseHandler
	termRepo repositories.TermRepository
}

func RegisterTermRoutes(tr repositories.TermRepository, mux *http.ServeMux) {
	h := &TermHandler{termRepo: tr}
	mux.HandleFunc("GET /api/terms/{id}", middleware.Use(h.GetTermDetail, middleware.Logger, middleware.CORS(nil)))
	mux.HandleFunc("GET /api/categories/{categoryId}/terms", middleware.Use(h.ListTermsByCategory, middleware.Logger, middleware.CORS(nil)))
}

func (h *TermHandler) GetTermDetail(w http.ResponseWriter, r *http.Request) {
	termID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "Invalid term ID", nil)
		return
	}

	term, err := h.termRepo.GetByID(r.Context(), termID)
	if err != nil {
		var serverErr *servererrors.ServerError
		if errors.As(err, &serverErr) {
			switch serverErr.Kind {
			case servererrors.KindNotFound:
				h.ResponseJSON(w, http.StatusNotFound, "Term not found", nil)
			default:
				h.ResponseJSON(w, http.StatusInternalServerError, "Internal server error", nil)
				log.Printf("TermHandler.GetTermDetail: failed to get term: %v", err)
			}
			return
		}

		h.ResponseJSON(w, http.StatusInternalServerError, "Internal server error", nil)
		log.Printf("TermHandler.GetTermDetail: failed to get term: %v", err)
		return
	}

	h.ResponseJSON(w, http.StatusOK, "Success", term)
}

// ListTermsByCategory 获取分类下的术语列表
func (h *TermHandler) ListTermsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := uuid.Parse(r.PathValue("categoryId"))
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	var req types.ListTermsRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 验证并设置limit
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	terms, hasMore, err := h.termRepo.ListByCategory(r.Context(), categoryID, req.LastID, req.Limit)
	if err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "Internal server error", nil)
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

	h.ResponseJSON(w, http.StatusOK, "Success", response)
}
