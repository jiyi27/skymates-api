package server

import (
	"net/http"
	"skymates-api/internal/repositories"
)

// Server 表示 HTTP 服务器
type Server struct {
	handler http.Handler
}

//// NewServer 创建并配置一个新的服务器实例
//func NewServer(repos *Repositories) *Server {
//	// 创建基础路由器
//	mux := http.NewServeMux()
//
//	// 创建服务层
//	services = service.NewServices()
//
//	// 注册 API v1 版本的路由
//	v1.RegisterRoutes(mux, services)
//
//	// 应用全局中间件
//	// 注意：中间件的应用顺序是从下到上，最后应用的中间件最先执行
//	var handler http.Handler = mux
//
//	return &Server{handler: handler}
//}

// Start 启动 HTTP 服务器
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}

// Repositories 包含所有仓库实例
type Repositories struct {
	UserRepository    repositories.UserRepository
	TermRepository    repositories.TermRepository
	PostRepository    repositories.PostRepository
	CommentRepository repositories.CommentRepository
}
