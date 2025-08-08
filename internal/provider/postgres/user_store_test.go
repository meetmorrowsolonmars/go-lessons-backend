//go:build integration && local

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain"
)

func TestUserStore_Create(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("Should create a new user", func(t *testing.T) {
		// Prepare.
		user := domain.User{
			Email:    "user_1@example.com",
			FullName: "Ivan Ivanov",
		}

		// Act.
		store := UserStore{
			db: db,
		}

		result, err := store.Create(ctx, user)

		// Check.
		assert.NoError(t, err)

		user.ID = result.ID
		user.CreateTime = result.CreateTime

		assert.Equal(t, user, result)
		assert.Greater(t, result.ID, int64(0))
		assert.NotEmpty(t, result.CreateTime)
	})
}
