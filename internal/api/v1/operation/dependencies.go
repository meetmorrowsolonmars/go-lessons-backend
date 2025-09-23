package operation

import (
	"context"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type OperationService interface {
	Create(ctx context.Context, operation model.Operation) (model.Operation, error)
	GetByAccountID(ctx context.Context, account int64, limit int64, offset int64) ([]model.Operation, error)
}

type JWTProvider interface {
	Validate(token string) (model.AuthClaims, error)
}
