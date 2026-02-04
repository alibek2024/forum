package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alibek2024/forum/internal/http/middleware"
	"github.com/alibek2024/forum/internal/models"
	"github.com/alibek2024/forum/internal/service"
)

type PostHandler struct {
	postService    *service.PostService
	commentService *service.CommentService
}

func NewPostHandler(p *service.PostService, c *service.CommentService) *PostHandler {
	return &PostHandler{postService: p, commentService: c}
}

// Главная страница с фильтрами
func (h *PostHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	user, _ := r.Context().Value(middleware.UserKey).(*models.User)

	// Собираем параметры из URL
	category := r.URL.Query().Get("category")
	filter := r.URL.Query().Get("filter")

	userID := 0
	if user != nil {
		userID = user.ID
	}

	// ВЫЗЫВАЕМ ОДИН УНИВЕРСАЛЬНЫЙ МЕТОД
	// Теперь сервис сам решит, как фильтровать, и при этом подтянет лайки
	posts, err := h.postService.GetAllPosts(category, filter, userID)

	if err != nil {
		log.Printf("ОШИБКА: %v", err)
		http.Error(w, "Ошибка загрузки постов", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "index", map[string]interface{}{
		"Posts": posts,
		"User":  user,
	})
}

// Создание поста (POST)
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// Парсим мультипарт (фото + текст)
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "Файлы слишком большие", http.StatusBadRequest)
			return
		}

		// 1. Получаем строку категорий и превращаем в слайс
		categoriesRaw := r.FormValue("categories")
		tags := strings.Split(categoriesRaw, ",")
		var cleanTags []string
		for _, t := range tags {
			t = strings.TrimSpace(t)
			if t != "" {
				cleanTags = append(cleanTags, strings.ToLower(t))
			}
		}

		// 2. Обработка файлов (как у тебя и было)
		files := r.MultipartForm.File["images"]
		var savedPaths []string
		var mainImagePath string

		for i, header := range files {
			file, err := header.Open()
			if err != nil {
				continue
			}

			fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + header.Filename
			path := "static/uploads/" + fileName
			dst, _ := os.Create("./web/" + path)
			io.Copy(dst, file)
			dst.Close()
			file.Close()

			if i == 0 {
				mainImagePath = path
			} else {
				savedPaths = append(savedPaths, path)
			}
		}

		// 3. Создаем объект поста
		post := &models.Post{
			UserID:   user.ID,
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
			ImageURL: mainImagePath,
		}

		// 4. Сохраняем всё в БД через твой новый метод Create
		postID, err := h.postService.CreatePost(post, cleanTags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 5. Добавляем остальные фото в post_images
		if len(savedPaths) > 0 {
			h.postService.AddImages(postID, savedPaths)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "create_post", map[string]interface{}{"User": user})
}

// UPDATE: Страница редактирования и сохранение изменений
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 1. ОБРАБОТКА СОХРАНЕНИЯ (POST)
	if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Файл слишком большой", http.StatusBadRequest)
			return
		}

		// Берем ID из скрытого поля формы <input type="hidden" name="id" ...>
		id, _ := strconv.Atoi(r.FormValue("id"))

		// Проверяем права: пост должен принадлежать пользователю
		oldPost, err := h.postService.GetPostByID(id)
		if err != nil || oldPost.UserID != user.ID {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		imagePath := oldPost.ImageURL // По умолчанию оставляем старое фото

		// Проверяем, загружено ли НОВОЕ фото
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + handler.Filename
			imagePath = "static/uploads/" + newFileName

			f, err := os.OpenFile("./web/"+imagePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err == nil {
				defer f.Close()
				io.Copy(f, file)
			}
		}

		// Создаем обновленный объект
		updatedPost := &models.Post{
			ID:       id,
			UserID:   user.ID,
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
			ImageURL: imagePath,
		}

		if err := h.postService.UpdatePost(updatedPost); err != nil {
			http.Error(w, "Ошибка обновления в базе", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post?id="+strconv.Itoa(id), http.StatusSeeOther)
		return
	}

	// 2. ОТОБРАЖЕНИЕ ФОРМЫ (GET)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	if post.UserID != user.ID {
		http.Error(w, "Вы не автор этого поста", http.StatusForbidden)
		return
	}

	renderTemplate(w, "edit_post", map[string]interface{}{
		"Post": post,
		"User": user,
	})
}

// DELETE: Удаление поста
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	// В сервисе должна быть проверка: post.UserID == user.ID
	if err := h.postService.DeletePost(postID, user.ID); err != nil {
		http.Error(w, "Удаление невозможно", http.StatusForbidden)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// READ (One): Просмотр конкретного поста и его комментариев
func (h *PostHandler) Show(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// ВАЖНО: Проверяем ошибку получения комментариев
	comments, err := h.commentService.GetCommentsByPostID(postID)
	if err != nil {
		// Если ошибка есть, логируем её, но можем вернуть пустой список, чтобы страница не падала
		comments = []models.Comment{}
	}

	// Получаем пользователя (он может быть nil, и это нормально для публичного просмотра)
	user, _ := r.Context().Value(middleware.UserKey).(*models.User)

	renderTemplate(w, "post_view", map[string]interface{}{
		"Post":     post,
		"Comments": comments,
		"User":     user,
	})
}
