package models

import "time"

type Session struct {
	ID        string // UUID
	UserID    int
	ExpiresAt time.Time
}
