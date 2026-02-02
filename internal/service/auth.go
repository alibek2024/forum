package service

import (
	"errors"
	"time"

	"github.com/alibek2024/forum/internal/models"
	"github.com/alibek2024/forum/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo    *repository.UserRepo
	SessionRepo *repository.SessionRepo
}

func NewAuthService(u *repository.UserRepo, s *repository.SessionRepo) *AuthService {
	return &AuthService{UserRepo: u, SessionRepo: s}
}

// Register user
func (s *AuthService) Register(email, username, password string) (*models.User, error) {
	existing, _ := s.UserRepo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := s.UserRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login user
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Create session
	sessionID := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	sess := &models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	if err := s.SessionRepo.Create(sess); err != nil {
		return "", err
	}

	return sessionID, nil
}
