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
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthService(u *repository.UserRepository, s *repository.SessionRepository) *AuthService {
	return &AuthService{userRepo: u, sessionRepo: s}
}

// Register — хеширует пароль и создает пользователя
func (s *AuthService) Register(email, username, password string) error {
	// Хеширование пароля (Бонусное задание ТЗ)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	return s.userRepo.CreateUser(user)
}

// Login — проверяет данные и генерирует UUID сессию
func (s *AuthService) Login(email, password string) (*models.Session, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// Сравнение хеша и пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// Создание сессии с UUID (Бонусное задание ТЗ)
	session := &models.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Сессия на 24 часа
	}

	if err := s.sessionRepo.Save(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) Logout(sessionID string) error {
	return s.sessionRepo.Delete(sessionID)
}

func (s *AuthService) GetUserBySession(sessionID string) (*models.User, error) {
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		s.sessionRepo.Delete(sessionID)
		return nil, errors.New("сессия истекла")
	}

	// Тут можно добавить метод в userRepo GetByID
	return &models.User{ID: session.UserID}, nil
}
