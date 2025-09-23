package user

import (
	"context"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type Store interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByID(ctx context.Context, id int64) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type AccountStore interface {
	CreateDefault(ctx context.Context, userID int64) (model.Account, error)
}
