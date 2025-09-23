package user

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

//go:embed testdata/create_request_success.json
var createRequestSuccess []byte

//go:embed testdata/create_request_invalid_body.xml
var createRequestInvalidBody []byte

func TestHandler_Create(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

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
				userService := NewUserServiceMock(mc)

				user := model.User{
					ID:         10,
					Email:      "user@example.com",
					Password:   "password123",
					FullName:   "Ivan Ivanov",
					CreateTime: now,
				}

				userService.CreateMock.
					When(minimock.AnyContext, model.User{
						Email:    user.Email,
						Password: user.Password,
						FullName: user.FullName,
					}).
					Then(user, nil)

				return &Handler{
					userService: userService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestSuccess,
			wantResp: json.RawMessage(`{ "id": 10 }`),
			wantCode: http.StatusCreated,
		},
		{
			name: "invalid request body",
			handler: func(mc *minimock.Controller) *Handler {
				userService := NewUserServiceMock(mc)

				return &Handler{
					userService: userService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestInvalidBody,
			wantResp: json.RawMessage(`{ "message": "Invalid request body: invalid character '<' looking for beginning of value" }`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "create user error",
			handler: func(mc *minimock.Controller) *Handler {
				userService := NewUserServiceMock(mc)

				user := model.User{
					ID:         10,
					Email:      "user@example.com",
					Password:   "password123",
					FullName:   "Ivan Ivanov",
					CreateTime: now,
				}

				userService.CreateMock.
					When(minimock.AnyContext, model.User{
						Email:    user.Email,
						Password: user.Password,
						FullName: user.FullName,
					}).
					Then(model.User{}, errors.New("some error"))

				return &Handler{
					userService: userService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestSuccess,
			wantResp: json.RawMessage(`{ "message": "Create user error" }`),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /v1/users", handler.Create)

			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/users", bytes.NewBuffer(tc.req))
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
