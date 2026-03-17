package models

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	IsActive     bool
	LastLogin    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
