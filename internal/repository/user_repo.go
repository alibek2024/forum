package repository

import (
	"database/sql"

	"github.com/alibek2024/forum/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	// Обязательный INSERT query
	query := `INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, user.Email, user.Username, user.PasswordHash)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	// Обязательный SELECT query
	query := `SELECT id, email, username, password_hash FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}
