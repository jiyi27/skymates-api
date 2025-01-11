package server

import (
	"net/http"
	"skymates-api/internal/handler"
	"skymates-api/internal/repositories"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(repos *Repositories) *Server {
	mux := http.NewServeMux()

	handler.RegisterUserRoutes(repos.UserRepository, mux)
	handler.RegisterTermRoutes(repos.TermRepository, mux)
	handler.RegisterPostRoutes(repos.PostRepository, mux)
	handler.RegisterCommentRoutes(repos.CommentRepository, mux)

	return &Server{mux: mux}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

type Repositories struct {
	UserRepository    repositories.UserRepository
	TermRepository    repositories.TermRepository
	PostRepository    repositories.PostRepository
	CommentRepository repositories.CommentRepository
}
