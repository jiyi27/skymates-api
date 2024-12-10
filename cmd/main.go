package main

import (
	"log"
	"skymates-api/internal/server"
	"skymates-api/internal/service"
	"skymates-api/internal/service/impl"
)

func initDatabase() service.Database {
	// 初始化数据库连接
	// 返回数据库实例
	return nil // 这里需要返回具体的数据库实现
}

func initUserService() service.UserService {
	db := initDatabase()
	return impl.NewUserService(db)
}

func initPostService() service.PostService {
	db := initDatabase()
	return impl.NewPostService(db)
}

func initCommentService() service.CommentService {
	db := initDatabase()
	return impl.NewCommentService(db)
}

func main() {
	// 初始化各个服务
	services := &server.Services{
		UserService:    initUserService(),
		PostService:    initPostService(),
		CommentService: initCommentService(),
	}

	// 创建并启动服务器
	srv := server.NewServer(services)
	log.Fatal(srv.Start(":8080"))
}
