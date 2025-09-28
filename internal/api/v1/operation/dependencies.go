package operation

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

//go:generate minimock -i OperationService

type OperationService interface {
	Create(ctx context.Context, operation model.Operation) (model.Operation, error)
	Update(ctx context.Context, id uuid.UUID, amount decimal.Decimal, categoryID int64, description string) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByAccountID(ctx context.Context, account int64, limit int64, offset int64) ([]model.Operation, error)
}

type JWTProvider interface {
	Validate(token string) (model.AuthClaims, error)
}
