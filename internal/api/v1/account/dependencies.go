package account

import (
	"context"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

//go:generate minimock -i AccountService

type AccountService interface {
	Create(ctx context.Context, account model.Account) (model.Account, error)
	GetByID(ctx context.Context, id int64) (model.Account, error)
	GetAccountsByUserID(ctx context.Context, userID int64) ([]model.Account, error)
}

type JWTProvider interface {
	Validate(token string) (model.AuthClaims, error)
}
