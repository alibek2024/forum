package repository

import (
	"github.com/alibek2024/forum/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	// Используем NamedExec, чтобы передать структуру целиком.
	// Названия после двоеточия (:email) должны совпадать с тегами db:"email" в модели.
	query := `INSERT INTO users (email, username, password_hash) VALUES (:email, :username, :password_hash)`
	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, username, password_hash FROM users WHERE email = ?`

	// Метод Get автоматически заменяет QueryRow + Scan
	err := r.db.Get(user, query, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, username, created_at FROM users WHERE id = ?`

	// Снова используем Get — sqlx сам разложит данные по полям структуры
	err := r.db.Get(user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET username = :username, email = :email WHERE id = :id`
	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
