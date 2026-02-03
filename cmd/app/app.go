package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alibek2024/forum/internal/db"
	"github.com/alibek2024/forum/internal/http/handlers"
	"github.com/alibek2024/forum/internal/http/middleware"
	"github.com/alibek2024/forum/internal/http/router"
	"github.com/alibek2024/forum/internal/repository"
	"github.com/alibek2024/forum/internal/service"
)

func main() {
	// 1. БД
	sqliteDB, err := db.InitSQLite("./forum.db", "./migrations/schema.sql")
	if err != nil {
		log.Fatal("Ошибка БД: ", err)
	}

	// 2. Репозитории
	userRepo := repository.NewUserRepository(sqliteDB)
	sessionRepo := repository.NewSessionRepository(sqliteDB)
	postRepo := repository.NewPostRepository(sqliteDB)
	commRepo := repository.NewCommentRepository(sqliteDB)
	likeRepo := repository.NewLikeRepository(sqliteDB)

	// 3. Сервисы
	authSvc := service.NewAuthService(userRepo, sessionRepo)
	postSvc := service.NewPostService(postRepo, likeRepo)
	commSvc := service.NewCommentService(commRepo, likeRepo)
	likeSvc := service.NewLikeService(likeRepo)

	// 4. Хендлеры
	authH := handlers.NewAuthHandler(authSvc)
	postH := handlers.NewPostHandler(postSvc, commSvc)
	commH := handlers.NewCommentHandler(commSvc)
	likeH := handlers.NewLikeHandler(likeSvc)

	// 5. Middleware
	authM := middleware.NewAuthMiddleware(authSvc)

	// 6. Роутер
	appRouter := router.NewRouter(authH, postH, commH, likeH, authM)

	fmt.Println("Форум запущен на http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", appRouter))
}
