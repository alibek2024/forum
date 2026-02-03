package service

import (
	"github.com/alibek2024/forum/internal/repository"
)

type LikeService struct {
	repo *repository.LikeRepository
}

func NewLikeService(r *repository.LikeRepository) *LikeService {
	return &LikeService{repo: r}
}

func (s *LikeService) HandleLike(userID, targetID int, targetType string, value int) error {
	// 1. Проверяем, есть ли уже такая реакция от этого пользователя
	// В репозитории мы использовали INSERT OR REPLACE, но для логики
	// "нажал второй раз — удалил" можно добавить проверку здесь или в репо.

	// Для простоты по ТЗ: просто устанавливаем значение (1 или -1)
	return s.repo.SetLike(userID, targetID, targetType, value)
}

func (s *LikeService) RemoveLike(userID, targetID int, targetType string) error {
	return s.repo.DeleteLike(userID, targetID, targetType)
}
