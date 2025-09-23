package auth

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

//go:embed testdata/login_request_success.json
var loginRequestSuccess []byte

//go:embed testdata/login_request_invalid_body.xml
var loginRequestInvalidBody []byte

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	token, _ := jwt.New(jwt.SigningMethodHS256).SignedString([]byte("secret"))

	cases := []struct {
		name     string
		handler  func(mc *minimock.Controller) *Handler
		req      []byte
		wantResp []byte
		wantCode int
	}{
		{
			name: "success",
			handler: func(mc *minimock.Controller) *Handler {
				authService := NewAuthServiceMock(mc)

				authService.LoginMock.
					When(minimock.AnyContext, "user@example.com", "password123").
					Then(token, nil)

				return &Handler{
					authService: authService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      loginRequestSuccess,
			wantResp: json.RawMessage(fmt.Sprintf(`{ "access_token": "%s" }`, token)),
			wantCode: http.StatusOK,
		},
		{
			name: "unauthorized",
			handler: func(mc *minimock.Controller) *Handler {
				authService := NewAuthServiceMock(mc)

				authService.LoginMock.
					When(minimock.AnyContext, "user@example.com", "password123").
					Then("", model.ErrNotAuthorized)

				return &Handler{
					authService: authService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      loginRequestSuccess,
			wantResp: json.RawMessage(`{ "message": "User not authorized" }`),
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "invalid request body",
			handler: func(mc *minimock.Controller) *Handler {
				authService := NewAuthServiceMock(mc)

				return &Handler{
					authService: authService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      loginRequestInvalidBody,
			wantResp: json.RawMessage(`{ "message": "Invalid request body: invalid character '<' looking for beginning of value" }`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "login error",
			handler: func(mc *minimock.Controller) *Handler {
				authService := NewAuthServiceMock(mc)

				authService.LoginMock.
					When(minimock.AnyContext, "user@example.com", "password123").
					Then("", errors.New("some error"))

				return &Handler{
					authService: authService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      loginRequestSuccess,
			wantResp: json.RawMessage(`{ "message": "Login error" }`),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /v1/login", handler.Login)

			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/login", bytes.NewBuffer(tc.req))
			require.NoError(t, err)

			resp, err := server.Client().Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			gotResp, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.wantCode, resp.StatusCode)
			assert.JSONEq(t, string(tc.wantResp), string(gotResp))
		})
	}
}
