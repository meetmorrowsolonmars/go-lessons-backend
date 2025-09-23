package user

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/examples/middleware/httpmiddleware"

	apimiddleware "github.com/meetmorrowsolonmars/education-pet-project/internal/api/middleware"
)

type Handler struct {
	userService UserService
	logger      *slog.Logger
}

func NewHandler(userService UserService, logger *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger.With(slog.String("logger", "user_handler")),
	}
}

func (h *Handler) Register(
	mux *http.ServeMux,
	promMiddleware httpmiddleware.Middleware,
	authMiddleware apimiddleware.AuthMiddleware,
) {
	mux.Handle("POST /v1/users",
		promMiddleware.WrapHandler("POST /v1/users", http.HandlerFunc(h.Create)))
	mux.Handle("GET /v1/users/{user_id}",
		promMiddleware.WrapHandler("GET /v1/users/{user_id}", authMiddleware.WrapHandler(http.HandlerFunc(h.GetByID))))
}
