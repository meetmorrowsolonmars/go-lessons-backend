package user

import (
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api/middleware"
)

type Handler struct {
	userService UserService
	jwtProvider JWTProvider
	logger      *slog.Logger
}

func NewHandler(userService UserService, jwtProvider JWTProvider, logger *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		jwtProvider: jwtProvider,
		logger:      logger.With(slog.String("logger", "user_handler")),
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/users", h.Create)
	mux.Handle("GET /v1/users/{user_id}", middleware.AuthMiddleware(h.jwtProvider)(http.HandlerFunc(h.GetByID)))
}
