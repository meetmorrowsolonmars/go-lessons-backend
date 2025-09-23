package user

import (
	"context"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByID(ctx context.Context, id int64) (model.User, error)
}

type JWTProvider interface {
	Validate(token string) (model.AuthClaims, error)
}
