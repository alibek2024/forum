package models

import "time"

type Session struct {
	ID        string    `db:"id"`
	UserID    int       `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
}
