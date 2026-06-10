package domain

import "context"

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
}

// Safe context key type
type contextKey string

const UserIDKey contextKey = "user_id"
