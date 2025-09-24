//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

func TestUserStore_Create(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "success",
			run: func(t *testing.T) {
				store := &UserStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Act.
				user := model.User{
					Email:    "user@example.com",
					Password: "password123",
					FullName: "Ivan Ivanov",
				}

				gotUser, gotErr := store.Create(ctx, user)

				// Check.
				assert.NoError(t, gotErr)

				row := db.QueryRow(ctx, getUserByIDQuery, gotUser.ID)

				var expected model.User
				err := row.Scan(
					&expected.ID,
					&expected.Email,
					&expected.Password,
					&expected.FullName,
					&expected.CreateTime,
				)
				assert.NoError(t, err)

				assert.Equal(t, expected, gotUser)

				// TODO: Check default account.

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

func TestUserStore_GetByID(t *testing.T) {
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
				store := &UserStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Prepare.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)

				// Act.
				gotUser, gotErr := store.GetByID(ctx, userID)

				gotUser.CreateTime = gotUser.CreateTime.UTC()

				// Check.
				assert.NoError(t, gotErr)

				expected := model.User{
					ID:         1,
					Email:      "empty@test.com",
					Password:   "",
					FullName:   "Ivan Ivanov",
					CreateTime: now,
				}
				assert.Equal(t, expected, gotUser)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users, accounts, operations;")
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				store := &UserStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Act.
				gotUser, gotErr := store.GetByID(ctx, userID)

				// Check.
				assert.ErrorIs(t, gotErr, model.ErrNotFound)
				assert.Zero(t, gotUser)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}

func TestUserStore_GetByEmail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)
	email := "empty@test.com"

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "success",
			run: func(t *testing.T) {
				store := &UserStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Prepare.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)

				// Act.
				gotUser, gotErr := store.GetByEmail(ctx, email)

				gotUser.CreateTime = gotUser.CreateTime.UTC()

				// Check.
				assert.NoError(t, gotErr)

				expected := model.User{
					ID:         1,
					Email:      "empty@test.com",
					Password:   "",
					FullName:   "Ivan Ivanov",
					CreateTime: now,
				}
				assert.Equal(t, expected, gotUser)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users, accounts, operations;")
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				store := &UserStore{
					db:  db,
					now: func() time.Time { return now },
				}

				// Act.
				gotUser, gotErr := store.GetByEmail(ctx, "nonexisting@example.com")

				// Check.
				assert.ErrorIs(t, gotErr, model.ErrNotFound)
				assert.Zero(t, gotUser)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}
