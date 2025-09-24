//go:build integration

package postgres

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

var (
	//go:embed testdata/create_test_operations.sql
	createTestOperationsQuery string
)

func TestOperationStore_Create(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := uuid.MustParse("52c99b00-cc95-42d9-94e6-cf24a411bcd9")
	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "success",
			run: func(t *testing.T) {
				store := &OperationStore{
					db:    db,
					newID: func() (uuid.UUID, error) { return id, nil },
					now:   func() time.Time { return now },
				}

				// Prepare.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)
				_, err = db.Exec(ctx, createTestAccountsQuery)
				assert.NoError(t, err)

				// Act.
				operation := model.Operation{
					UserID:      1,
					AccountID:   1,
					Type:        model.OperationTypeDebit,
					Amount:      decimal.NewFromFloat(100.55),
					Description: "Salary",
				}

				gotOperation, gotErr := store.Create(ctx, operation)

				// Check.
				assert.NoError(t, gotErr)

				row := db.QueryRow(ctx, getOperationByIdQuery, gotOperation.ID)

				var expected model.Operation
				err = row.Scan(
					&expected.ID,
					&expected.UserID,
					&expected.AccountID,
					&expected.Type,
					&expected.Amount,
					&expected.Description,
					&expected.CreateTime,
				)
				assert.NoError(t, err)

				expected.CreateTime = expected.CreateTime.UTC()

				assert.Equal(t, expected, gotOperation)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users, accounts, operations;")
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}

func TestOperationStore_Update(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := uuid.MustParse("52c99b00-cc95-42d9-94e6-cf24a411bcd9")
	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "success",
			run: func(t *testing.T) {
				store := &OperationStore{
					db:    db,
					newID: func() (uuid.UUID, error) { return id, nil },
					now:   func() time.Time { return now },
				}

				// Prepare.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)
				_, err = db.Exec(ctx, createTestAccountsQuery)
				assert.NoError(t, err)
				_, err = db.Exec(ctx, createTestOperationsQuery)
				assert.NoError(t, err)

				// Act.
				amount := decimal.NewFromFloat(300.55)
				description := "Crypto"

				gotErr := store.Update(ctx, id, amount, description)

				// Check.
				assert.NoError(t, gotErr)

				row := db.QueryRow(ctx, getOperationByIdQuery, id)

				var gotOperation model.Operation
				err = row.Scan(
					&gotOperation.ID,
					&gotOperation.UserID,
					&gotOperation.AccountID,
					&gotOperation.Type,
					&gotOperation.Amount,
					&gotOperation.Description,
					&gotOperation.CreateTime,
				)
				assert.NoError(t, err)

				assert.Equal(t, amount, gotOperation.Amount)
				assert.Equal(t, description, gotOperation.Description)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users, accounts, operations;")
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				store := &OperationStore{
					db:    db,
					newID: func() (uuid.UUID, error) { return id, nil },
					now:   func() time.Time { return now },
				}

				// Act.
				amount := decimal.NewFromFloat(300.55)
				description := "Crypto"

				gotErr := store.Update(ctx, id, amount, description)

				// Check.
				assert.ErrorIs(t, gotErr, model.ErrNotFound)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}
