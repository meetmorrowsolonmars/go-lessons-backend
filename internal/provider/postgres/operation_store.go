package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

var (
	//go:embed queries/create_operation.sql
	createOperationQuery string

	//go:embed queries/update_operation.sql
	updateOperationQuery string

	//go:embed queries/delete_operation.sql
	deleteOperationQuery string

	//go:embed queries/get_operation_by_id.sql
	getOperationByIdQuery string

	//go:embed queries/get_operations_by_account_id.sql
	getOperationsByAccountIdQuery string
)

type OperationStore struct {
	db    *pgxpool.Pool
	newID func() (uuid.UUID, error)
	now   func() time.Time
}

func NewOperationStore(db *pgxpool.Pool) *OperationStore {
	return &OperationStore{
		db:    db,
		newID: func() (uuid.UUID, error) { return uuid.NewV7() },
		now:   func() time.Time { return time.Now().UTC() },
	}
}

func (s *OperationStore) Create(ctx context.Context, operation model.Operation) (model.Operation, error) {
	id, err := s.newID()
	if err != nil {
		return model.Operation{}, fmt.Errorf("generate operation id: %w", err)
	}

	operation.ID = id
	operation.CreateTime = s.now()
	categoryID := sql.NullInt64{
		Int64: operation.CategoryID,
		Valid: operation.CategoryID != 0,
	}

	_, err = s.db.Exec(
		ctx,
		createOperationQuery,
		operation.ID,
		operation.UserID,
		operation.AccountID,
		operation.Type,
		categoryID,
		operation.Amount,
		operation.Description,
		operation.CreateTime,
	)
	if err != nil {
		return model.Operation{}, fmt.Errorf("create operation: %w", err)
	}

	return operation, nil
}

func (s *OperationStore) Update(
	ctx context.Context,
	id uuid.UUID,
	amount decimal.Decimal,
	categoryID int64,
	description string,
) error {
	sqlCategoryID := sql.NullInt64{
		Int64: categoryID,
		Valid: categoryID != 0,
	}

	result, err := s.db.Exec(ctx, updateOperationQuery, id, amount, description, sqlCategoryID)
	if err != nil {
		return fmt.Errorf("update operation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("operation not found: %w", model.ErrNotFound)
	}

	return nil
}

func (s *OperationStore) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.Exec(ctx, deleteOperationQuery, id)
	if err != nil {
		return fmt.Errorf("delete operation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("operation not found: %w", model.ErrNotFound)
	}

	return nil
}

func (s *OperationStore) GetByAccountID(
	ctx context.Context,
	accountID int64,
	limit int64,
	offset int64,
) ([]model.Operation, error) {
	rows, err := s.db.Query(ctx, getOperationsByAccountIdQuery, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get operations by account id: %w", err)
	}

	defer rows.Close()

	operations := make([]model.Operation, 0, limit)

	for rows.Next() {
		var operation model.Operation
		var categoryID sql.NullInt64

		err = rows.Scan(
			&operation.ID,
			&operation.UserID,
			&operation.AccountID,
			&operation.Type,
			&categoryID,
			&operation.Amount,
			&operation.Description,
			&operation.CreateTime,
		)
		if err != nil {
			return nil, fmt.Errorf("get operations by account id: %w", err)
		}

		operation.CategoryID = categoryID.Int64

		operations = append(operations, operation)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get operations by account id: %w", err)
	}

	return operations, nil
}
