package auth

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/examples/middleware/httpmiddleware"

	apimiddleware "github.com/meetmorrowsolonmars/education-pet-project/internal/api/middleware"
)

type Handler struct {
	authService AuthService
	logger      *slog.Logger
}

func NewHandler(
	authService AuthService,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger.With(slog.String("logger", "auth_handler")),
	}
}

func (h *Handler) Register(
	mux *http.ServeMux,
	promMiddleware httpmiddleware.Middleware,
	_ apimiddleware.AuthMiddleware,
) {
	mux.HandleFunc("POST /v1/login", promMiddleware.WrapHandler("POST /v1/login", http.HandlerFunc(h.Login)))
}
