package operation

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
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

//go:embed testdata/get_account_operations_response_success.json
var getUserAccountsResponseSuccess []byte

func TestHandler_GetAccountOperations(t *testing.T) {
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
				operationService := NewOperationServiceMock(mc)

				operations := []model.Operation{
					{
						ID:          uuid.MustParse("da270184-3a9e-4082-9509-c48a6f8c505f"),
						UserID:      10,
						AccountID:   10,
						Type:        model.OperationTypeDebit,
						Amount:      decimal.NewFromFloat(160.5),
						Description: "Salary",
						CreateTime:  now,
					},
					{
						ID:          uuid.MustParse("516f29da-d014-4594-a28e-771ca71b80c2"),
						UserID:      10,
						AccountID:   10,
						Type:        model.OperationTypeCredit,
						Amount:      decimal.NewFromFloat(160.5),
						Description: "New iPhone",
						CreateTime:  now,
					},
				}
				operationService.GetByAccountIDMock.
					When(minimock.AnyContext, 10, 20, 0).
					Then(operations, nil)

				return &Handler{
					operationService: operationService,
					logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/operations/10",
			wantResp: getUserAccountsResponseSuccess,
			wantCode: http.StatusOK,
		},
		{
			name: "invalid limit",
			handler: func(mc *minimock.Controller) *Handler {
				operationService := NewOperationServiceMock(mc)

				return &Handler{
					operationService: operationService,
					logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/operations/10?limit=invalid",
			wantResp: json.RawMessage(`{ "message": "Invalid request body: strconv.ParseInt: parsing \"invalid\": invalid syntax" }`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid offset",
			handler: func(mc *minimock.Controller) *Handler {
				operationService := NewOperationServiceMock(mc)

				return &Handler{
					operationService: operationService,
					logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/operations/10?offset=invalid",
			wantResp: json.RawMessage(`{ "message": "Invalid request body: strconv.ParseInt: parsing \"invalid\": invalid syntax" }`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid account id",
			handler: func(mc *minimock.Controller) *Handler {
				operationService := NewOperationServiceMock(mc)

				return &Handler{
					operationService: operationService,
					logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/operations/invalid_account_id",
			wantResp: json.RawMessage(`{ "message": "Invalid request body: strconv.ParseInt: parsing \"invalid_account_id\": invalid syntax" }`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "get account operations error",
			handler: func(mc *minimock.Controller) *Handler {
				operationService := NewOperationServiceMock(mc)

				operationService.GetByAccountIDMock.
					When(minimock.AnyContext, 10, 20, 0).
					Then(nil, errors.New("some error"))

				return &Handler{
					operationService: operationService,
					logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			path:     "/v1/operations/10",
			wantResp: json.RawMessage(`{ "message": "Get account operations error" }`),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /v1/operations/{account_id}", handler.GetAccountOperations)

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
