package main

import (
	"database/sql"
	"log"
	"net/http"
	v1 "skymates-api/api/v1"
	"skymates-api/internal/repository"
	"skymates-api/internal/service"
	"skymates-api/pkg/middleware"
)

func main() {
	// 1. 初始化数据库连接
	db, err := repository.NewMySQLDatabase()
	if err != nil {
		log.Fatal("init database failed: ", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("close database failed: ", err)
		}
	}(db)

	// 2. 初始化仓库
	userRepository := repository.NewUserRepository(db)
	services := &service.Services{UserService: service.NewUserService(userRepository)}

	// 3. 创建 HTTP 路由
	router := http.NewServeMux()
	v1.RegisterRoutes(router, services)

	// 4. 添加中间件
	handler := addGlobalMiddlewares(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func addGlobalMiddlewares(handler http.Handler) http.Handler {
	// 先应用日志中间件, 记录所有请求
	handler = middleware.Logger(handler)
	handler = middleware.CORS(handler)
	return handler
}
