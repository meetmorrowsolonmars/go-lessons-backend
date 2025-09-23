package operation

import (
	"context"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type Store interface {
	Create(ctx context.Context, operation model.Operation) (model.Operation, error)
	GetByAccountID(ctx context.Context, accountID int64, limit int64, offset int64) ([]model.Operation, error)
}

type UserStore interface {
	GetByID(ctx context.Context, id int64) (model.User, error)
}

type AccountStore interface {
	GetByID(ctx context.Context, id int64) (model.Account, error)
}
