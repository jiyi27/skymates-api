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

	// 初始化各个 handler
	userHandler := handler.NewUserHandler(repos.UserRepository)
	postHandler := handler.NewPostHandler(repos.PostRepository)
	commentHandler := handler.NewCommentHandler(repos.CommentRepository)

	// 注册路由
	userHandler.RegisterRoutes(mux)
	postHandler.RegisterRoutes(mux)
	commentHandler.RegisterRoutes(mux)

	return &Server{mux: mux}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

type Repositories struct {
	UserRepository    repositories.UserRepository
	PostRepository    repositories.PostRepository
	CommentRepository repositories.CommentRepository
}
