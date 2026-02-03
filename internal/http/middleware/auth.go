package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/alibek2024/forum/internal/service"
)

// Определяем тип для ключа в контексте, чтобы избежать коллизий
type contextKey string

const UserKey contextKey = "user"

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(s *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: s}
}

func (m *AuthMiddleware) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Пытаемся достать куку
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// Если куки нет, пользователь гость.
			// Просто идем дальше, хендлеры сами решат, что делать.
			next.ServeHTTP(w, r)
			return
		}

		// 2. Проверяем сессию через сервис
		user, err := m.authService.GetUserBySession(cookie.Value)
		if err != nil {
			// Если сессия протухла или неверна — удаляем куку
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now().Add(-1 * time.Hour),
				Path:    "/",
			})
			next.ServeHTTP(w, r)
			return
		}

		// 3. Добавляем пользователя в контекст запроса
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
