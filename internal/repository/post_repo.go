package repository

import (
	"database/sql"

	"github.com/alibek2024/forum/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post, categories []string) error {
	tx, err := r.db.Begin() // Используем транзакцию, так как пишем в две таблицы
	if err != nil {
		return err
	}

	// 1. Создаем пост
	res, err := tx.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		post.UserID, post.Title, post.Content)
	if err != nil {
		tx.Rollback()
		return err
	}

	postID, _ := res.LastInsertId()

	// 2. Привязываем категории
	for _, catName := range categories {
		// Сначала находим ID категории или создаем её (хотя категории обычно предопределены)
		var catID int
		err := tx.QueryRow("SELECT id FROM categories WHERE name = ?", catName).Scan(&catID)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, catID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at 
		FROM posts p 
		JOIN users u ON p.user_id = u.id 
		ORDER BY p.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.UserID, &p.AuthorName, &p.Title, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
