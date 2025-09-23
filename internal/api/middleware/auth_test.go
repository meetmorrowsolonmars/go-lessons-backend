package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	token, _ := jwt.New(jwt.SigningMethodHS256).SignedString([]byte("secret"))

	cases := []struct {
		name     string
		handler  func(mc *minimock.Controller) http.Handler
		header   string
		wantCode int
	}{
		{
			name: "successful authorization",
			handler: func(mc *minimock.Controller) http.Handler {
				jwtProvider := NewJWTProviderMock(mc)

				expectedClaims := model.AuthClaims{
					UserID: 100,
				}

				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					claims, _ := model.AuthClaimsValue(r.Context())
					assert.Equal(mc, expectedClaims, claims)

					_, _ = fmt.Fprintln(w, "Hello, world!")
				})

				jwtProvider.ValidateMock.
					When(token).
					Then(expectedClaims, nil)

				return AuthMiddleware(jwtProvider)(handler)
			},
			header:   fmt.Sprintf("Bearer %s", token),
			wantCode: http.StatusOK,
		},
		{
			name: "authorization header missing",
			handler: func(mc *minimock.Controller) http.Handler {
				jwtProvider := NewJWTProviderMock(mc)

				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = fmt.Fprintln(w, "Hello, world!")
				})

				return AuthMiddleware(jwtProvider)(handler)
			},
			header:   "",
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "authorization header format is invalid",
			handler: func(mc *minimock.Controller) http.Handler {
				jwtProvider := NewJWTProviderMock(mc)

				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = fmt.Fprintln(w, "Hello, world!")
				})

				return AuthMiddleware(jwtProvider)(handler)
			},
			header:   fmt.Sprintf("Invalid %s", token),
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "validate token error",
			handler: func(mc *minimock.Controller) http.Handler {
				jwtProvider := NewJWTProviderMock(mc)

				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = fmt.Fprintln(w, "Hello, world!")
				})

				jwtProvider.ValidateMock.
					When(token).
					Then(model.AuthClaims{}, errors.New("some error"))

				return AuthMiddleware(jwtProvider)(handler)
			},
			header:   fmt.Sprintf("Bearer %s", token),
			wantCode: http.StatusUnauthorized,
		},
	}

	// TODO: Should refactor this test to follow the same pattern as other tests in this layer?

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			req := httptest.NewRequest(http.MethodGet, "/v1/accounts", nil)
			resp := httptest.NewRecorder()

			if tc.header != "" {
				req.Header.Add("Authorization", tc.header)
			}

			handler := tc.handler(mc)

			handler.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
		})
	}
}
