package handler

import (
	"net/http"
	"skymates-api/internal/service"
)

type PostHandler struct {
	BaseHandler
	postService service.PostService
}

func NewPostHandler(ps service.PostService) *PostHandler {
	return &PostHandler{postService: ps}
}

func (h *PostHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/posts", h.handlePosts)
	mux.HandleFunc("/api/posts/{id}", h.handlePost)
}

func (h *PostHandler) handlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreatePost(w, r)
	case http.MethodGet:
		h.ListPosts(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *PostHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetPost(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// 实现创建帖子逻辑
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	// 实现获取帖子逻辑
}

func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	// 实现获取帖子列表逻辑
}
