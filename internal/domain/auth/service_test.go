package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

func TestService_Login(t *testing.T) {
	t.Parallel()

	token, _ := jwt.New(jwt.SigningMethodHS256).SignedString([]byte("secret"))

	cases := []struct {
		name      string
		service   func(mc *minimock.Controller) *Service
		email     string
		password  string
		wantToken string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			service: func(mc *minimock.Controller) *Service {
				userService := NewUserServiceMock(mc)
				jwtProvider := NewJWTProviderMock(mc)

				user := model.User{
					Password: "password123",
				}
				err := user.HashPassword()
				require.NoError(t, err)

				userService.GetByEmailMock.
					When(minimock.AnyContext, "user@example.com").
					Then(user, nil)

				jwtProvider.GenerateMock.
					When(user).
					Then(token, nil)

				return &Service{
					userService: userService,
					jwtProvider: jwtProvider,
				}
			},
			email:     "user@example.com",
			password:  "password123",
			wantToken: token,
			wantErr:   assert.NoError,
		},
		{
			name: "user not found",
			service: func(mc *minimock.Controller) *Service {
				userService := NewUserServiceMock(mc)
				jwtProvider := NewJWTProviderMock(mc)

				userService.GetByEmailMock.
					When(minimock.AnyContext, "user@example.com").
					Then(model.User{}, model.ErrNotFound)

				return &Service{
					userService: userService,
					jwtProvider: jwtProvider,
				}
			},
			email:     "user@example.com",
			password:  "password123",
			wantToken: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, model.ErrNotAuthorized, i...)
			},
		},
		{
			name: "get user error",
			service: func(mc *minimock.Controller) *Service {
				userService := NewUserServiceMock(mc)
				jwtProvider := NewJWTProviderMock(mc)

				userService.GetByEmailMock.
					When(minimock.AnyContext, "user@example.com").
					Then(model.User{}, errors.New("some error"))

				return &Service{
					userService: userService,
					jwtProvider: jwtProvider,
				}
			},
			email:     "user@example.com",
			password:  "password123",
			wantToken: "",
			wantErr:   assert.Error,
		},
		{
			name: "invalid password",
			service: func(mc *minimock.Controller) *Service {
				userService := NewUserServiceMock(mc)
				jwtProvider := NewJWTProviderMock(mc)

				user := model.User{
					Password: "password123",
				}
				err := user.HashPassword()
				require.NoError(t, err)

				userService.GetByEmailMock.
					When(minimock.AnyContext, "user@example.com").
					Then(user, nil)

				return &Service{
					userService: userService,
					jwtProvider: jwtProvider,
				}
			},
			email:     "user@example.com",
			password:  "password123456789",
			wantToken: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, model.ErrNotAuthorized, i...)
			},
		},
		{
			name: "generate token error",
			service: func(mc *minimock.Controller) *Service {
				userService := NewUserServiceMock(mc)
				jwtProvider := NewJWTProviderMock(mc)

				user := model.User{
					Password: "password123",
				}
				err := user.HashPassword()
				require.NoError(t, err)

				userService.GetByEmailMock.
					When(minimock.AnyContext, "user@example.com").
					Then(user, nil)

				jwtProvider.GenerateMock.
					When(user).
					Then("", errors.New("some error"))

				return &Service{
					userService: userService,
					jwtProvider: jwtProvider,
				}
			},
			email:     "user@example.com",
			password:  "password123",
			wantToken: "",
			wantErr:   assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			service := tc.service(mc)

			gotToken, gotErr := service.Login(context.Background(), tc.email, tc.password)

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantToken, gotToken)
		})
	}
}
