package auth

import (
	"log/slog"
	"net/http"
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

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/login", h.Login)
}
