package user

import (
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

//go:embed testdata/get_by_id_response_success.json
var getByIdResponseSuccess []byte

func TestHandler_GetByID(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		handler  func(mc *minimock.Controller) *Handler
		path     string
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

				userService.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(user, nil)

				return &Handler{
					userService: userService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/users/10",
			wantResp: getByIdResponseSuccess,
			wantCode: http.StatusOK,
		},
		{
			name: "invalid user id",
			handler: func(mc *minimock.Controller) *Handler {
				userService := NewUserServiceMock(mc)

				return &Handler{
					userService: userService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/users/invalid_user_id",
			wantResp: json.RawMessage(`{ "message": "Invalid request body: strconv.ParseInt: parsing \"invalid_user_id\": invalid syntax" }`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "user not found",
			handler: func(mc *minimock.Controller) *Handler {
				userService := NewUserServiceMock(mc)

				userService.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.User{}, model.ErrNotFound)

				return &Handler{
					userService: userService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/users/10",
			wantResp: json.RawMessage(`{ "message": "User not found" }`),
			wantCode: http.StatusNotFound,
		},
		{
			name: "get account by id error",
			handler: func(mc *minimock.Controller) *Handler {
				userService := NewUserServiceMock(mc)

				userService.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.User{}, errors.New("some error"))

				return &Handler{
					userService: userService,
					logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/users/10",
			wantResp: json.RawMessage(`{ "message": "Get user by id error" }`),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /v1/users/{user_id}", handler.GetByID)

			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL+tc.path, nil)
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
