package repository

import (
	"strings"

	"github.com/alibek2024/forum/internal/models"
	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Базовый запрос, который мы будем переиспользовать везде.
// Он сразу считает лайки/дизлайки и подтягивает имя автора.
const basePostSelect = `
    SELECT 
        p.id, p.user_id, p.title, p.content, 
        COALESCE(p.image_url, '') AS image_url, 
        p.created_at, u.username as author_name
    FROM posts p
    JOIN users u ON p.user_id = u.id`

// CREATE: Создание поста и привязка категорий
// Create создает пост и привязывает категории в одной транзакции
func (r *PostRepository) Create(post *models.Post, categories []string) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 1. Вставляем сам пост
	query := `INSERT INTO posts (user_id, title, content, image_url) 
              VALUES (:user_id, :title, :content, :image_url)`

	res, err := tx.NamedExec(query, post)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	postID := int(lastID)

	// 2. Привязываем категории, используя наш внутренний метод
	if err := r.attachCategories(tx, postID, categories); err != nil {
		return 0, err
	}

	// 3. Фиксируем изменения
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return postID, nil
}

// AddCategoriesToPost пригодится для добавления тегов к уже созданному посту (например, при Edit)
func (r *PostRepository) AddCategoriesToPost(postID int, categoriesStr string) error {
	// Превращаем строку в слайс чистых тегов
	tags := strings.Split(categoriesStr, ",")
	var cleanTags []string
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag != "" {
			cleanTags = append(cleanTags, tag)
		}
	}

	// Открываем транзакцию специально для этого действия
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.attachCategories(tx, postID, cleanTags); err != nil {
		return err
	}

	return tx.Commit()
}

// Приватный метод-помощник, чтобы не дублировать код.
// Он работает с транзакцией, которую ему передали.
func (r *PostRepository) attachCategories(tx *sqlx.Tx, postID int, tags []string) error {
	for _, tag := range tags {
		// Создаем категорию, если её нет
		_, err := tx.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", tag)
		if err != nil {
			return err
		}

		// Получаем ID категории
		var catID int
		err = tx.Get(&catID, "SELECT id FROM categories WHERE name = ?", tag)
		if err != nil {
			return err
		}

		// Связываем пост с категорией
		_, err = tx.Exec("INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, catID)
		if err != nil {
			return err
		}
	}
	return nil
}

// READ (ALL): Все посты
func (r *PostRepository) GetAll(category string, filter string, userID int) ([]models.Post, error) {
	query := basePostSelect + " WHERE 1=1 "
	var args []interface{}

	if category != "" {
		query += " AND p.id IN (SELECT post_id FROM post_categories pc JOIN categories c ON pc.category_id = c.id WHERE c.name = ?) "
		args = append(args, category)
	}

	if filter == "created" && userID != 0 {
		query += " AND p.user_id = ? "
		args = append(args, userID)
	}

	if filter == "liked" && userID != 0 {
		// ВНИМАНИЕ: используем value = 1, так как в твоей схеме колонка называется value
		query += " AND p.id IN (SELECT target_id FROM likes WHERE user_id = ? AND target_type = 'post' AND value = 1) "
		args = append(args, userID)
	}

	query += " GROUP BY p.id ORDER BY p.created_at DESC"
	return r.fetchPosts(query, args...)
}

func (r *PostRepository) GetCategoriesByPostID(postID int) ([]string, error) {
	var categories []string
	query := `
        SELECT c.name 
        FROM categories c 
        JOIN post_categories pc ON c.id = pc.category_id 
        WHERE pc.post_id = ?`

	err := r.db.Select(&categories, query, postID)
	return categories, err
}

// READ (ONE): Один пост по ID
func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	var post models.Post

	// 1. Получаем основной пост + имя автора + категории одной строкой
	query := `
        SELECT 
            p.*, 
            u.username as author_name, 
            GROUP_CONCAT(DISTINCT c.name) as categories_str
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN post_categories pc ON p.id = pc.post_id
        LEFT JOIN categories c ON pc.category_id = c.id
        WHERE p.id = ?
        GROUP BY p.id`

	err := r.db.Get(&post, query, id)
	if err != nil {
		return nil, err
	}

	// 2. Превращаем строку категорий "tag1,tag2" в слайс []string
	if post.CategoriesStr != "" {
		post.Categories = strings.Split(post.CategoriesStr, ",")
	}

	// 3. Подтягиваем все картинки (твой существующий код)
	var images []string
	err = r.db.Select(&images, "SELECT path FROM post_images WHERE post_id = ?", id)
	if err == nil {
		post.Images = images
	}

	return &post, nil
}

func (r *PostRepository) AddImages(postID int, paths []string) error {
	for _, path := range paths {
		_, err := r.db.Exec("INSERT INTO post_images (post_id, path) VALUES (?, ?)", postID, path)
		if err != nil {
			return err
		}
	}
	return nil
}

// UPDATE: Изменение контента
func (r *PostRepository) Update(post *models.Post) error {
	query := `
        UPDATE posts 
        SET title = ?, content = ?, image_url = ? 
        WHERE id = ? AND user_id = ?`

	_, err := r.db.Exec(query, post.Title, post.Content, post.ImageURL, post.ID, post.UserID)
	return err
}

// DELETE: Удаление со всеми связями
func (r *PostRepository) Delete(postID, userID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// --- ФИЛЬТРЫ (Теперь с лайками!) ---

func (r *PostRepository) GetByCategory(categoryName string) ([]models.Post, error) {
	query := basePostSelect + `
       JOIN post_categories pc ON p.id = pc.post_id
       JOIN categories c ON pc.category_id = c.id
       WHERE c.name = ?
       ORDER BY p.created_at DESC`
	return r.fetchPosts(query, categoryName)
}

func (r *PostRepository) GetByUserID(userID int) ([]models.Post, error) {
	query := basePostSelect + " WHERE p.user_id = ? ORDER BY p.created_at DESC"
	return r.fetchPosts(query, userID)
}

func (r *PostRepository) GetLikedByUser(userID int) ([]models.Post, error) {
	query := basePostSelect + `
       JOIN likes l ON p.id = l.target_id
       WHERE l.user_id = ? AND l.target_type = 'post' AND l.value = 1
       ORDER BY p.created_at DESC`
	return r.fetchPosts(query, userID)
}

// Вспомогательный метод для выполнения запросов
func (r *PostRepository) fetchPosts(query string, args ...interface{}) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Select(&posts, query, args...)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		// 1. Инициализация слайсов
		posts[i].Images = []string{}
		posts[i].Categories = []string{}

		// 2. Подтягиваем картинки
		imgQuery := `SELECT path FROM post_images WHERE post_id = ?`
		r.db.Select(&posts[i].Images, imgQuery, posts[i].ID)

		// 3. ПОДТЯГИВАЕМ КАТЕГОРИИ (Reddit-style)
		// Без этого твой {{range .Categories}} в шаблоне будет пустым!
		catQuery := `
          SELECT c.name 
          FROM categories c 
          JOIN post_categories pc ON c.id = pc.category_id 
          WHERE pc.post_id = ?`
		r.db.Select(&posts[i].Categories, catQuery, posts[i].ID)
	}

	return posts, nil
}

func (r *PostRepository) GetPostByID(id int) (*models.Post, error) {
	var post models.Post
	query := `
        SELECT p.*, u.username as author_name, 
               GROUP_CONCAT(c.name) AS categories_str 
        FROM posts p
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN post_categories pc ON p.id = pc.post_id
        LEFT JOIN categories c ON pc.category_id = c.id
        WHERE p.id = ?
        GROUP BY p.id`

	err := r.db.Get(&post, query, id)
	return &post, err
}
