package account

import (
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api/middleware"
)

type Handler struct {
	accountService AccountService
	jwtProvider    JWTProvider
	logger         *slog.Logger
}

func NewHandler(accountService AccountService, jwtProvider JWTProvider, logger *slog.Logger) *Handler {
	return &Handler{
		accountService: accountService,
		jwtProvider:    jwtProvider,
		logger:         logger.With(slog.String("logger", "account_handler")),
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.Handle("POST /v1/accounts", middleware.AuthMiddleware(h.jwtProvider)(http.HandlerFunc(h.Create)))
	mux.Handle("GET /v1/accounts", middleware.AuthMiddleware(h.jwtProvider)(http.HandlerFunc(h.GetUserAccounts)))
	mux.Handle("GET /v1/accounts/{account_id}", middleware.AuthMiddleware(h.jwtProvider)(http.HandlerFunc(h.GetByID)))
}
