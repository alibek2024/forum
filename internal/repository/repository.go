package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User    *UserRepository
	Session *SessionRepository
	Post    *PostRepository
	Comment *CommentRepository
	Like    *LikeRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Session: NewSessionRepository(db),
		Post:    NewPostRepository(db),
		Comment: NewCommentRepository(db),
		Like:    NewLikeRepository(db),
	}
}
