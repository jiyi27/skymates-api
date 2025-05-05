package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
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

	// 2. 将 *sql.DB 转换为 *sqlx.DB
	// *sql.DB 和 *sqlx.DB 共享同一个底层连接池, 调用 db.Close() 会关闭整个连接池，因此只需关闭一次即可
	sqlxDB := sqlx.NewDb(db, "mysql")

	// 3. 初始化仓库
	userRepository := repository.NewUserRepository(sqlxDB)
	termRepository := repository.NewTermRepository(sqlxDB)

	// 4. 初始化服务
	services := &service.Services{
		UserService: service.NewUserService(userRepository),
		TermService: service.NewTermService(termRepository),
	}

	// 5. 创建 HTTP 路由
	router := http.NewServeMux()
	v1.RegisterRoutes(router, services)

	// 6. 添加中间件
	handler := addGlobalMiddlewares(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func addGlobalMiddlewares(handler http.Handler) http.Handler {
	// 先应用日志中间件, 记录所有请求
	handler = middleware.Logger(handler)
	handler = middleware.CORS(handler)
	return handler
}
