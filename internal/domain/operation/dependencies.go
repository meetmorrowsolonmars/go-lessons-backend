package operation

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

//go:generate minimock -i Store

type Store interface {
	Create(ctx context.Context, operation model.Operation) (model.Operation, error)
	Update(ctx context.Context, id uuid.UUID, amount decimal.Decimal, categoryID int64, description string) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByAccountID(ctx context.Context, accountID int64, limit int64, offset int64) ([]model.Operation, error)
}

//go:generate minimock -i UserStore

type UserStore interface {
	GetByID(ctx context.Context, id int64) (model.User, error)
}

//go:generate minimock -i AccountStore

type AccountStore interface {
	GetByID(ctx context.Context, id int64) (model.Account, error)
}

//go:generate minimock -i CategoryStore

type CategoryStore interface {
	GetByID(ctx context.Context, id int64) (model.Category, error)
}
