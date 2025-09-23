package auth

import (
	"context"
)

//go:generate minimock -i AuthService

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
}
