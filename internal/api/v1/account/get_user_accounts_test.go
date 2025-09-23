package account

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

//go:embed testdata/get_user_accounts_response_success.json
var getUserAccountsResponseSuccess []byte

func TestHandler_GetUserAccounts(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)
	userID := int64(10)

	cases := []struct {
		name     string
		handler  func(mc *minimock.Controller) *Handler
		wantResp []byte
		wantCode int
	}{
		{
			name: "success",
			handler: func(mc *minimock.Controller) *Handler {
				accountService := NewAccountServiceMock(mc)

				accounts := []model.Account{
					{
						ID:         1,
						UserID:     userID,
						Title:      "default",
						IsDefault:  true,
						CreateTime: now,
					},
					{
						ID:         2,
						UserID:     userID,
						Title:      "piggy bank",
						IsDefault:  false,
						CreateTime: now,
					},
				}

				accountService.GetAccountsByUserIDMock.
					When(minimock.AnyContext, userID).
					Then(accounts, nil)

				return &Handler{
					accountService: accountService,
					// TODO: Create noop logger and redirect it to os.DevNull.
					logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			wantResp: getUserAccountsResponseSuccess,
			wantCode: http.StatusOK,
		},
		{
			name: "get account error",
			handler: func(mc *minimock.Controller) *Handler {
				accountService := NewAccountServiceMock(mc)

				accountService.GetAccountsByUserIDMock.
					When(minimock.AnyContext, userID).
					Then([]model.Account{}, errors.New("some error"))

				return &Handler{
					accountService: accountService,
					// TODO: Create noop logger and redirect it to os.DevNull.
					logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			wantResp: json.RawMessage(`{ "message": "Get user accounts error" }`),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /v1/accounts", func(w http.ResponseWriter, r *http.Request) {
				ctx := model.WithAuthClaims(r.Context(), model.AuthClaims{
					UserID: userID,
				})
				r = r.WithContext(ctx)
				handler.GetUserAccounts(w, r)
			})

			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL+"/v1/accounts", nil)
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
