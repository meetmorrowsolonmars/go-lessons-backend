package operation

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
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	userID := int64(10)

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
				operationService := NewOperationServiceMock(mc)

				operation := model.Operation{
					ID:          uuid.MustParse("52c99b00-cc95-42d9-94e6-cf24a411bcd9"),
					UserID:      userID,
					AccountID:   10,
					Type:        model.OperationTypeDebit,
					Amount:      decimal.NewFromFloat(160.5),
					Description: "Salary",
					CreateTime:  now,
				}

				operationService.CreateMock.
					When(minimock.AnyContext, model.Operation{
						UserID:      operation.UserID,
						AccountID:   operation.AccountID,
						Type:        operation.Type,
						Amount:      operation.Amount,
						Description: operation.Description,
					}).
					Then(operation, nil)

				return &Handler{
					operationService: operationService,
					logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestSuccess,
			wantResp: json.RawMessage(`{ "id": "52c99b00-cc95-42d9-94e6-cf24a411bcd9" }`),
			wantCode: http.StatusCreated,
		},
		{
			name: "invalid request body",
			handler: func(mc *minimock.Controller) *Handler {
				operationService := NewOperationServiceMock(mc)

				return &Handler{
					operationService: operationService,
					logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestInvalidBody,
			wantResp: json.RawMessage(`{ "message": "Invalid request body: invalid character '<' looking for beginning of value" }`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "create operation error",
			handler: func(mc *minimock.Controller) *Handler {
				operationService := NewOperationServiceMock(mc)

				operation := model.Operation{
					ID:          uuid.MustParse("52c99b00-cc95-42d9-94e6-cf24a411bcd9"),
					UserID:      userID,
					AccountID:   10,
					Type:        model.OperationTypeDebit,
					Amount:      decimal.NewFromFloat(160.5),
					Description: "Salary",
					CreateTime:  now,
				}

				operationService.CreateMock.
					When(minimock.AnyContext, model.Operation{
						UserID:      operation.UserID,
						AccountID:   operation.AccountID,
						Type:        operation.Type,
						Amount:      operation.Amount,
						Description: operation.Description,
					}).
					Then(model.Operation{}, errors.New("some error"))

				return &Handler{
					operationService: operationService,
					logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
				}
			},
			req:      createRequestSuccess,
			wantResp: json.RawMessage(`{ "message": "Create operation error" }`),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /v1/operations", func(w http.ResponseWriter, r *http.Request) {
				ctx := model.WithAuthClaims(r.Context(), model.AuthClaims{
					UserID: userID,
				})
				r = r.WithContext(ctx)
				handler.Create(w, r)
			})

			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/operations", bytes.NewBuffer(tc.req))
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
