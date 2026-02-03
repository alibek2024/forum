package handlers

import (
	"net/http"
	"strconv"

	"github.com/alibek2024/forum/internal/http/middleware"
	"github.com/alibek2024/forum/internal/models"
	"github.com/alibek2024/forum/internal/service"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(s *service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: s}
}

func (h *LikeHandler) HandleLike(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok {
		http.Error(w, "You must be logged in to like", http.StatusUnauthorized)
		return
	}

	var targetID int
	var targetType string
	var value int

	if r.Method == http.MethodPost {
		// Если пришло из формы (для постов)
		targetID, _ = strconv.Atoi(r.FormValue("target_id"))
		targetType = r.FormValue("target_type")
		value, _ = strconv.Atoi(r.FormValue("value"))
	} else {
		// Если пришло из ссылки (для комментариев)
		targetID, _ = strconv.Atoi(r.URL.Query().Get("id"))
		targetType = r.URL.Query().Get("type")
		value, _ = strconv.Atoi(r.URL.Query().Get("value"))
	}

	// ВАЖНО: Добавь проверку для отладки, если опять упадет
	if targetType == "" {
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		return
	}

	err := h.likeService.HandleLike(user.ID, targetID, targetType, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
