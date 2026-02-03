package models

import "time"

type Comment struct {
	ID         int       `db:"id"`
	PostID     int       `db:"post_id"`
	UserID     int       `db:"user_id"`
	AuthorName string    `db:"author_name"` // Появится после JOIN users
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"` // Время создания
	Likes      int       `db:"likes"`
	Dislikes   int       `db:"dislikes"`
}
