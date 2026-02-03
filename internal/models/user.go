package models

import "time"

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"` // Важно: как в схеме!
	CreatedAt    time.Time `db:"created_at"`
}
