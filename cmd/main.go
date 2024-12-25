package main

import (
	"log"
	"skymates-api/internal/repositories/impl"
	"skymates-api/internal/server"
)

func main() {
	db, err := impl.NewPostgresDB()
	if err != nil {
		log.Fatal("init database failed: ", err)
	}
	defer db.Close()

	repos := &server.Repositories{
		UserRepository:    impl.NewPostgresUserRepository(db),
		PostRepository:    impl.NewPostgresPostRepository(db),
		CommentRepository: impl.NewPostgresCommentRepository(db),
	}

	srv := server.NewServer(repos)
	log.Fatal(srv.Start(":8080"))
}
