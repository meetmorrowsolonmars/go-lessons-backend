package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeRequest[LoginRequest](r)
	if err != nil {
		h.logger.Error("Decode login request", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if errors.Is(err, model.ErrNotAuthorized) {
		h.logger.Error("User not authorized",
			slog.String("email", req.Email),
			slog.String("error", err.Error()),
		)

		api.EncodeErrorf(w, http.StatusUnauthorized, "User not authorized")

		return
	}
	if err != nil {
		h.logger.Error("Login error", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Login error: %s", err)

		return
	}

	api.EncodeSuccess(w, http.StatusOK, LoginResponse{
		AccessToken: token,
	})
}
