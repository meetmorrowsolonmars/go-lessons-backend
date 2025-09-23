package account

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/examples/middleware/httpmiddleware"

	apimiddleware "github.com/meetmorrowsolonmars/education-pet-project/internal/api/middleware"
)

type Handler struct {
	accountService AccountService
	logger         *slog.Logger
}

func NewHandler(accountService AccountService, logger *slog.Logger) *Handler {
	return &Handler{
		accountService: accountService,
		logger:         logger.With(slog.String("logger", "account_handler")),
	}
}

func (h *Handler) Register(
	mux *http.ServeMux,
	promMiddleware httpmiddleware.Middleware,
	authMiddleware apimiddleware.AuthMiddleware,
) {
	mux.Handle("POST /v1/accounts",
		promMiddleware.WrapHandler("POST /v1/accounts",
			authMiddleware.WrapHandler(http.HandlerFunc(h.Create))))
	mux.Handle("GET /v1/accounts",
		promMiddleware.WrapHandler("GET /v1/accounts",
			authMiddleware.WrapHandler(http.HandlerFunc(h.GetUserAccounts))))
	mux.Handle("GET /v1/accounts/{account_id}",
		promMiddleware.WrapHandler("GET /v1/accounts/{account_id}",
			authMiddleware.WrapHandler(http.HandlerFunc(h.GetByID))))
}
