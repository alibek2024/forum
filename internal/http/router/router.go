package router

import (
	"net/http"

	"github.com/alibek2024/forum/internal/http/handlers"
	"github.com/alibek2024/forum/internal/http/middleware"
)

func NewRouter(
	authH *handlers.AuthHandler,
	postH *handlers.PostHandler,
	commH *handlers.CommentHandler,
	likeH *handlers.LikeHandler,
	authM *middleware.AuthMiddleware,
) *http.ServeMux {
	mux := http.NewServeMux()

	// СТАТИКА (CSS, JS)
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// AUTH
	mux.HandleFunc("/register", authH.Register)
	mux.HandleFunc("/login", authH.Login)
	mux.HandleFunc("/logout", authH.Logout)

	// POSTS (Публичные)
	mux.HandleFunc("/", authM.CheckAuth(postH.Index)) // Index сам решит, что показать гостю
	mux.HandleFunc("/post", authM.CheckAuth(postH.Show))

	// POSTS (Приватные - только для авторизованных)
	mux.HandleFunc("/post/create", authM.CheckAuth(postH.Create))
	mux.HandleFunc("/post/edit", authM.CheckAuth(postH.Update))
	mux.HandleFunc("/post/delete", authM.CheckAuth(postH.Delete))

	// COMMENTS
	mux.HandleFunc("/comment/create", authM.CheckAuth(commH.Create))
	mux.HandleFunc("/comment/edit", authM.CheckAuth(commH.Update))
	mux.HandleFunc("/comment/delete", authM.CheckAuth(commH.Delete))

	// LIKES
	mux.HandleFunc("/like", authM.CheckAuth(likeH.HandleLike))

	return mux
}
