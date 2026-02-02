package models

import "time"

type Comment struct {
	ID         int
	PostID     int
	UserID     int
	AuthorName string
	Content    string
	CreatedAt  time.Time
	Likes      int
	Dislikes   int
}
