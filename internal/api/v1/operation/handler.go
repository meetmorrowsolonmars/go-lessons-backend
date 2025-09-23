package operation

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/examples/middleware/httpmiddleware"

	apimiddleware "github.com/meetmorrowsolonmars/education-pet-project/internal/api/middleware"
)

type Handler struct {
	operationService OperationService
	logger           *slog.Logger
}

func NewHandler(
	operationService OperationService,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		operationService: operationService,
		logger:           logger.With(slog.String("logger", "operation_handler")),
	}
}

func (h *Handler) Register(
	mux *http.ServeMux,
	promMiddleware httpmiddleware.Middleware,
	authMiddleware apimiddleware.AuthMiddleware,
) {
	mux.Handle("POST /v1/operations",
		promMiddleware.WrapHandler("POST /v1/operations",
			authMiddleware.WrapHandler(http.HandlerFunc(h.Create))))
	mux.Handle("GET /v1/operations/{account_id}",
		promMiddleware.WrapHandler("GET /v1/operations/{account_id}",
			authMiddleware.WrapHandler(http.HandlerFunc(h.GetAccountOperations))),
	)
}
