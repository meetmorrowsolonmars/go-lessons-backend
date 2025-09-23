package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		service  func(mc *minimock.Controller) *Service
		user     model.User
		wantUser model.User
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)

				store.CreateMock.
					Set(func(ctx context.Context, user model.User) (model.User, error) {
						assert.Equal(mc, "user@example.com", user.Email)
						assert.Equal(mc, "Ivan Ivanov", user.FullName)
						assert.True(mc, user.CheckPassword("password123"))

						return model.User{
							ID:         1,
							Email:      "user@example.com",
							Password:   user.Password,
							FullName:   "Ivan Ivanov",
							CreateTime: now,
						}, nil
					})

				return &Service{
					store: store,
				}
			},
			user: model.User{
				Email:    "user@example.com",
				Password: "password123",
				FullName: "Ivan Ivanov",
			},
			wantUser: model.User{
				ID:         1,
				Email:      "user@example.com",
				FullName:   "Ivan Ivanov",
				CreateTime: now,
			},
			wantErr: assert.NoError,
		},
		{
			name: "create user error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)

				store.CreateMock.
					Set(func(ctx context.Context, user model.User) (model.User, error) {
						assert.Equal(mc, "user@example.com", user.Email)
						assert.Equal(mc, "Ivan Ivanov", user.FullName)
						assert.True(mc, user.CheckPassword("password123"))

						return model.User{}, errors.New("some error")
					})

				return &Service{
					store: store,
				}
			},
			user: model.User{
				Email:    "user@example.com",
				Password: "password123",
				FullName: "Ivan Ivanov",
			},
			wantUser: model.User{},
			wantErr:  assert.Error,
		},
		{
			name: "create account error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)

				store.CreateMock.
					Set(func(ctx context.Context, user model.User) (model.User, error) {
						assert.Equal(mc, "user@example.com", user.Email)
						assert.Equal(mc, "Ivan Ivanov", user.FullName)
						assert.True(mc, user.CheckPassword("password123"))

						return model.User{
							ID:         1,
							Email:      "user@example.com",
							Password:   user.Password,
							FullName:   "Ivan Ivanov",
							CreateTime: now,
						}, nil
					})

				return &Service{
					store: store,
				}
			},
			user: model.User{
				Email:    "user@example.com",
				Password: "password123",
				FullName: "Ivan Ivanov",
			},
			wantUser: model.User{},
			wantErr:  assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			service := tc.service(mc)

			gotUser, gotErr := service.Create(context.Background(), tc.user)

			tc.wantUser.Password = gotUser.Password

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantUser, gotUser)
		})
	}
}
