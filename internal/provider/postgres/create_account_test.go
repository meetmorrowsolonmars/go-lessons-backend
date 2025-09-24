//go:build integration

package postgres

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

var (
	//go:embed testdata/create_test_user.sql
	createTestUserQuery string

	//go:embed testdata/create_test_accounts.sql
	createTestAccountsQuery string
)

func TestAccountStore_Create(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)
	userID := int64(1)

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "success",
			run: func(t *testing.T) {
				store := &AccountStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Prepare.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)

				// Act.
				account := model.Account{
					UserID:    userID,
					Title:     "Piggy Bank",
					IsDefault: false,
				}

				gotAccount, gotErr := store.Create(ctx, account)

				// Check.
				assert.NoError(t, gotErr)

				row := db.QueryRow(ctx, getAccountByIDQuery, gotAccount.ID)

				var expected model.Account
				err = row.Scan(
					&expected.ID,
					&expected.UserID,
					&expected.Title,
					&expected.IsDefault,
					&expected.CreateTime,
				)
				assert.NoError(t, err)

				assert.Equal(t, expected, gotAccount)

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

func TestAccountStore_GetByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)
	accountID := int64(1)

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "success",
			run: func(t *testing.T) {
				store := &AccountStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Prepare.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)
				_, err = db.Exec(ctx, createTestAccountsQuery)
				assert.NoError(t, err)

				// Act.
				gotAccount, gotErr := store.GetByID(ctx, accountID)

				gotAccount.CreateTime = gotAccount.CreateTime.UTC()

				// Check.
				assert.NoError(t, gotErr)

				expected := model.Account{
					ID:         1,
					UserID:     1,
					Title:      "default",
					IsDefault:  true,
					CreateTime: now,
				}

				assert.Equal(t, expected, gotAccount)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users, accounts, operations;")
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				store := &AccountStore{
					db:  db,
					now: func() time.Time { return now },
				}

				accountID := int64(101)

				// Act.
				gotAccount, gotErr := store.GetByID(ctx, accountID)

				// Check.
				assert.ErrorIs(t, gotErr, model.ErrNotFound)
				assert.Zero(t, gotAccount)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}

func TestAccountStore_GetAccountsByUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)
	userID := int64(1)

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "success",
			run: func(t *testing.T) {
				store := &AccountStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Prepare.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)
				_, err = db.Exec(ctx, createTestAccountsQuery)
				assert.NoError(t, err)

				// Act.
				gotAccounts, gotErr := store.GetAccountsByUserID(ctx, userID)

				for i := range gotAccounts {
					gotAccounts[i].CreateTime = gotAccounts[i].CreateTime.UTC()
				}

				// Check.
				assert.NoError(t, gotErr)

				expected := []model.Account{
					{
						ID:         1,
						UserID:     1,
						Title:      "default",
						IsDefault:  true,
						CreateTime: now,
					},
					{
						ID:         2,
						UserID:     1,
						Title:      "Piggy Bank",
						IsDefault:  false,
						CreateTime: now,
					},
				}

				assert.Equal(t, expected, gotAccounts)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users, accounts, operations;")
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				store := &AccountStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Act.
				gotAccounts, gotErr := store.GetAccountsByUserID(ctx, userID)

				// Check.
				assert.NoError(t, gotErr)
				assert.Len(t, gotAccounts, 0)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}
