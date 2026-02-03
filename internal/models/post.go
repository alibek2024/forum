package models

import "time"

type Post struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	AuthorName string    `db:"author_name"`
	Title      string    `db:"title"`
	Content    string    `db:"content"`
	ImageURL   string    `db:"image_url"`
	Images     []string  `db:"-"`
	CreatedAt  time.Time `db:"created_at"`
	Likes      int       `db:"likes"`
	Dislikes   int       `db:"dislikes"`
	// Добавь эти поля ниже:
	CategoriesStr string   `db:"categories_str"` // Сюда SQL запишет "go,docker"
	Categories    []string `db:"-"`              // Это мы выведем в HTML
}
