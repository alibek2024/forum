package router

import (
	"net/http"

	"github.com/alibek2024/forum/internal/http/handlers"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	postHandler *handlers.PostHandler,
	commentHandler *handlers.CommentHandler,
	likeHandler *handlers.LikeHandler,
) http.Handler {

	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("forum is running"))
	})

	// auth
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/logout", authHandler.Logout)

	// posts
	mux.HandleFunc("/posts", postHandler.List)          // GET
	mux.HandleFunc("/posts/create", postHandler.Create) // POST

	// comments
	mux.HandleFunc("/comments/create", commentHandler.Create)

	// likes
	mux.HandleFunc("/likes", likeHandler.Toggle)

	return mux
}
