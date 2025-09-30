package operation

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/metric"
)

type Service struct {
	store         Store
	userStore     UserStore
	accountStore  AccountStore
	categoryStore CategoryStore
}

func NewService(
	store Store,
	userStore UserStore,
	accountStore AccountStore,
	categoryStore CategoryStore,
) *Service {
	return &Service{
		store:         store,
		userStore:     userStore,
		accountStore:  accountStore,
		categoryStore: categoryStore,
	}
}

func (s *Service) Create(ctx context.Context, operation model.Operation) (model.Operation, error) {
	// TODO: Decide where validate amount in api or domain layer. Amount always should be positive number.
	if operation.CategoryID != 0 {
		_, err := s.categoryStore.GetByID(ctx, operation.CategoryID)
		if err != nil {
			return model.Operation{}, fmt.Errorf("get category by id: %w", err)
		}
	}

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

func (s *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	amount decimal.Decimal,
	categoryID int64,
	description string,
) error {
	err := s.store.Update(ctx, id, amount, categoryID, description)
	if err != nil {
		return fmt.Errorf("update operation: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.store.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete operation: %w", err)
	}

	return nil
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
