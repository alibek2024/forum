package handlers

import (
	"net/http"
	"strconv"

	"github.com/alibek2024/forum/internal/http/middleware"
	"github.com/alibek2024/forum/internal/models"
	"github.com/alibek2024/forum/internal/service"
)

type CommentHandler struct {
	commentService *service.CommentService
}

// Передаем сервис при создании, чтобы хендлер мог им пользоваться
func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: s}
}

// Create — теперь рабочий метод
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем авторизацию через контекст (из Middleware)
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok {
		http.Error(w, "You must be logged in to comment", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  user.ID,
		Content: content,
	}

	if err := h.commentService.CreateComment(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем пользователя туда, откуда он пришел (на страницу поста)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// Update и Delete оставляем как у тебя, они написаны верно,
// если в Service реализованы проверки авторства.

// Update — редактирование комментария
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok || r.Method != http.MethodPost {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
	content := r.FormValue("content")

	if err := h.commentService.UpdateComment(commentID, user.ID, content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// Delete — удаление комментария
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if err := h.commentService.DeleteComment(commentID, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
