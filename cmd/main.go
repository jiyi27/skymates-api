package main

import (
	"log"
	"net/http"
	v1 "skymates-api/api/v1"
	"skymates-api/pkg/middleware"
)

func main() {
	//db, err := impl.NewPostgresDB()
	//if err != nil {
	//	log.Fatal("init database failed: ", err)
	//}
	//defer db.Close()
	//
	//userRepo := impl.NewPostgresUserRepository(db)
	//services := &service.Services{UserService: service.NewUserService(userRepo)}
	
	router := http.NewServeMux()
	v1.RegisterRoutes(router, nil)

	handler := addGlobalMiddlewares(router)
	log.Fatal(http.ListenAndServe(":8080", handler))

	//log.Print("Starting server...")
	//db, err := impl.NewPostgresDB()
	//if err != nil {
	//	log.Fatal("init database failed: ", err)
	//}
	//defer db.Close()
	//
	//repos := &server.Repositories{
	//	UserRepository:    impl.NewPostgresUserRepository(db),
	//	TermRepository:    impl.NewPostgresTermRepository(db),
	//	PostRepository:    impl.NewPostgresPostRepository(db),
	//	CommentRepository: impl.NewPostgresCommentRepository(db),
	//}
	//
	//srv := server.NewServer(repos)
	//log.Print("Server started at :8080")
	//log.Fatal(srv.Start(":8080"))
}

func addGlobalMiddlewares(handler http.Handler) http.Handler {
	// 先应用日志中间件, 记录所有请求
	handler = middleware.Logger(handler)
	handler = middleware.CORS(handler)
	return handler
}
