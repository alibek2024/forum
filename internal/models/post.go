package models

import "time"

type Post struct {
	ID         int
	UserID     int
	AuthorName string
	Title      string
	Content    string
	Categories []string
	CreatedAt  time.Time
	Likes      int
	Dislikes   int
}

// Вспомогательная структура для обработки лайков
type Vote struct {
	UserID     int
	TargetID   int
	TargetType string // "post" или "comment"
	Value      int    // 1 или -1
}
