package service

import (
	"errors"
	"strings"

	"github.com/alibek2024/forum/internal/models"
	"github.com/alibek2024/forum/internal/repository"
)

type PostService struct {
	repo     *repository.PostRepository
	likeRepo *repository.LikeRepository
}

func NewPostService(r *repository.PostRepository, l *repository.LikeRepository) *PostService {
	return &PostService{repo: r, likeRepo: l}
}

func (s *PostService) CreatePost(post *models.Post, categories []string) (int, error) {
	return s.repo.Create(post, categories)
}

func (s *PostService) GetAllPosts(category string, filter string, userID int) ([]models.Post, error) {
	posts, err := s.repo.GetAll(category, filter, userID)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		// 1. Лайки (через LikeRepo)
		l, d, _ := s.likeRepo.GetCounts(posts[i].ID, "post")
		posts[i].Likes = l
		posts[i].Dislikes = d

		// 2. Категории (если fetchPosts их вдруг не подтянул)
		if len(posts[i].Categories) == 0 {
			cats, _ := s.repo.GetCategoriesByPostID(posts[i].ID)
			posts[i].Categories = cats
		}
	}
	return posts, nil
}

func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	post, err := s.repo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	// 1. Превращаем строку "go,docker" в массив []string{"go", "docker"}
	if post.CategoriesStr != "" {
		post.Categories = strings.Split(post.CategoriesStr, ",")
	}

	// 2. Подтягиваем лайки (твой существующий код)
	likes, dislikes, _ := s.likeRepo.GetCounts(post.ID, "post")
	post.Likes = likes
	post.Dislikes = dislikes

	return post, nil
}

func (s *PostService) UpdatePost(post *models.Post) error {
	// Проверяем, существует ли пост и принадлежит ли он юзеру
	existing, err := s.repo.GetPostByID(post.ID)
	if err != nil {
		return errors.New("post not found")
	}
	if existing.UserID != post.UserID {
		return errors.New("you are not the author of this post")
	}

	return s.repo.Update(post)
}

func (s *PostService) DeletePost(postID, userID int) error {
	// Проверка прав
	existing, err := s.repo.GetPostByID(postID)
	if err != nil {
		return errors.New("post not found")
	}
	if existing.UserID != userID {
		return errors.New("permission denied")
	}

	return s.repo.Delete(postID, userID)
}

// Методы для фильтрации (по ТЗ)
func (s *PostService) GetPostsByUser(userID int) ([]models.Post, error) {
	return s.repo.GetByUserID(userID)
}

func (s *PostService) GetLikedPosts(userID int) ([]models.Post, error) {
	return s.repo.GetLikedByUser(userID)
}

func (s *PostService) GetPostsByCategory(category string) ([]models.Post, error) {
	return s.repo.GetByCategory(category)
}

func (s *PostService) AddImages(postID int, paths []string) error {
	return s.repo.AddImages(postID, paths)
}
