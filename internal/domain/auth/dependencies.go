package auth

import (
	"context"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type UserService interface {
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type JWTProvider interface {
	Generate(user model.User) (string, error)
}
