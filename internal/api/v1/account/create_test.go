package account

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
	userID := int64(100)

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
				accountService := NewAccountServiceMock(mc)

				account := model.Account{
					ID:         10,
					UserID:     userID,
					Title:      "piggy bank",
					IsDefault:  false,
					CreateTime: now,
				}

				accountService.CreateMock.
					When(minimock.AnyContext, model.Account{
						UserID: account.UserID,
						Title:  account.Title,
					}).
					Then(account, nil)

				return &Handler{
					accountService: accountService,
					// TODO: Create noop logger and redirect it to os.DevNull.
					logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestSuccess,
			wantResp: json.RawMessage(`{ "id": 10 }`),
			wantCode: http.StatusCreated,
		},
		{
			name: "invalid request body",
			handler: func(mc *minimock.Controller) *Handler {
				accountService := NewAccountServiceMock(mc)

				return &Handler{
					accountService: accountService,
					logger:         slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestInvalidBody,
			wantResp: json.RawMessage(`{ "message": "Invalid request body: invalid character '<' looking for beginning of value" }`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "create account error",
			handler: func(mc *minimock.Controller) *Handler {
				accountService := NewAccountServiceMock(mc)

				account := model.Account{
					UserID: userID,
					Title:  "piggy bank",
				}

				accountService.CreateMock.
					When(minimock.AnyContext, account).
					Then(model.Account{}, errors.New("some error"))

				return &Handler{
					accountService: accountService,
					logger:         slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestSuccess,
			wantResp: json.RawMessage(`{ "message": "Create account error" }`),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /v1/accounts", func(w http.ResponseWriter, r *http.Request) {
				ctx := model.WithAuthClaims(r.Context(), model.AuthClaims{
					UserID: userID,
				})
				handler.Create(w, r.WithContext(ctx))
			})

			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/accounts", bytes.NewBuffer(tc.req))
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
