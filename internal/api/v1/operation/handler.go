package operation

import (
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api/middleware"
)

type Handler struct {
	operationService OperationService
	jwtProvider      JWTProvider
	logger           *slog.Logger
}

func NewHandler(
	operationService OperationService,
	jwtProvider JWTProvider,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		operationService: operationService,
		jwtProvider:      jwtProvider,
		logger:           logger.With(slog.String("logger", "operation_handler")),
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.Handle("POST /v1/operations", middleware.AuthMiddleware(h.jwtProvider)(http.HandlerFunc(h.Create)))
	mux.Handle("GET /v1/operations/{account_id}", middleware.AuthMiddleware(h.jwtProvider)(http.HandlerFunc(h.GetAccountOperations)))
}
