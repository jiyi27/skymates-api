package handler

import (
	"net/http"
	"skymates-api/internal/repositories"
)

type CommentHandler struct {
	BaseHandler
	commentService repositories.CommentRepository
}

func NewCommentHandler(cs repositories.CommentRepository) *CommentHandler {
	return &CommentHandler{commentService: cs}
}

func (h *CommentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/posts/{postID}/comments", h.handleComments)
}

func (h *CommentHandler) handleComments(w http.ResponseWriter, r *http.Request) {
	//switch r.Method {
	//case http.MethodPost:
	//	h.CreateComment(w, r)
	//case http.MethodGet:
	//	h.ListComments(w, r)
	//default:
	//	http.NotFound(w, r)
	//}
}
