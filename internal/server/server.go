package server

import (
	"net/http"
	"skymates-api/internal/handler"
	"skymates-api/internal/service"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(services *Services) *Server {
	mux := http.NewServeMux()

	// 初始化各个 handler
	userHandler := handler.NewUserHandler(services.UserService)
	postHandler := handler.NewPostHandler(services.PostService)
	commentHandler := handler.NewCommentHandler(services.CommentService)

	// 注册路由
	userHandler.RegisterRoutes(mux)
	postHandler.RegisterRoutes(mux)
	commentHandler.RegisterRoutes(mux)

	return &Server{mux: mux}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

type Services struct {
	UserService    service.UserService
	PostService    service.PostService
	CommentService service.CommentService
}
