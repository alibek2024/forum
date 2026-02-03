package service

import "github.com/alibek2024/forum/internal/repository"

type Service struct {
	Auth    *AuthService
	Post    *PostService
	Comment *CommentService
	Like    *LikeService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth:    NewAuthService(repos.User, repos.Session),
		Post:    NewPostService(repos.Post, repos.Like),
		Comment: NewCommentService(repos.Comment, repos.Like),
		Like:    NewLikeService(repos.Like),
	}
}
