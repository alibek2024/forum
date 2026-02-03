package repository

import (
	"github.com/alibek2024/forum/internal/models"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Save(s *models.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, s.ID, s.UserID, s.ExpiresAt)
	return err
}

func (r *SessionRepository) GetByID(id string) (*models.Session, error) {
	query := `SELECT id, user_id, expires_at FROM sessions WHERE id = ?`
	row := r.db.QueryRow(query, id)

	s := &models.Session{}
	err := row.Scan(&s.ID, &s.UserID, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SessionRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}
