package auth

import (
	"context"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

//go:generate minimock -i UserService

type UserService interface {
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

//go:generate minimock -i JWTProvider

type JWTProvider interface {
	Generate(user model.User) (string, error)
}
