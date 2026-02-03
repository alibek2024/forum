package service

import (
	"errors"

	"github.com/alibek2024/forum/internal/models"
	"github.com/alibek2024/forum/internal/repository"
)

type CommentService struct {
	repo     *repository.CommentRepository
	likeRepo *repository.LikeRepository
}

func NewCommentService(r *repository.CommentRepository, l *repository.LikeRepository) *CommentService {
	return &CommentService{
		repo:     r,
		likeRepo: l,
	}
}

func (s *CommentService) CreateComment(comment *models.Comment) error {
	if comment.Content == "" {
		return errors.New("комментарий не может быть пустым")
	}
	return s.repo.Create(comment)
}

func (s *CommentService) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	comments, err := s.repo.GetByPostID(postID)
	if err != nil {
		return nil, err
	}

	// Подтягиваем лайки для каждого комментария
	for i := range comments {
		likes, dislikes, _ := s.likeRepo.GetCounts(comments[i].ID, "comment")
		comments[i].Likes = likes
		comments[i].Dislikes = dislikes
	}

	return comments, nil
}

func (s *CommentService) UpdateComment(commentID, userID int, content string) error {
	if content == "" {
		return errors.New("comment content cannot be empty")
	}

	// Проверяем авторство
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}
	if comment.UserID != userID {
		return errors.New("you can only edit your own comments")
	}

	return s.repo.Update(commentID, userID, content)
}

func (s *CommentService) DeleteComment(commentID, userID int) error {
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	// Только автор комментария (или автор поста, если захочешь) может удалить
	if comment.UserID != userID {
		return errors.New("not authorized to delete this comment")
	}

	return s.repo.Delete(commentID, userID)
}
