package operation

import (
	"context"
	"fmt"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/metric"
)

type Service struct {
	store        Store
	userStore    UserStore
	accountStore AccountStore
}

func NewService(
	store Store,
	userStore UserStore,
	accountStore AccountStore,
) *Service {
	return &Service{
		store:        store,
		userStore:    userStore,
		accountStore: accountStore,
	}
}

func (s *Service) Create(ctx context.Context, operation model.Operation) (model.Operation, error) {
	// TODO: Decide where validate amount in api or domain layer. Amount always should be positive number.
	_, err := s.userStore.GetByID(ctx, operation.UserID)
	if err != nil {
		return model.Operation{}, fmt.Errorf("get user by id: %w", err)
	}

	_, err = s.accountStore.GetByID(ctx, operation.AccountID)
	if err != nil {
		return model.Operation{}, fmt.Errorf("get account by id: %w", err)
	}

	operation, err = s.store.Create(ctx, operation)
	if err != nil {
		return model.Operation{}, fmt.Errorf("create operation: %w", err)
	}

	metric.OperationsTotalCounterInc(operation.Type)

	return operation, nil
}

func (s *Service) GetByAccountID(
	ctx context.Context,
	accountID int64,
	limit int64,
	offset int64,
) ([]model.Operation, error) {
	_, err := s.accountStore.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get account by id: %w", err)
	}

	operations, err := s.store.GetByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get operations: %w", err)
	}

	return operations, nil
}
