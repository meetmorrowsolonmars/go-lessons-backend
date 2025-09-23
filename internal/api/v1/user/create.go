package user

import (
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type CreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type CreateResponse struct {
	ID int64 `json:"id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeRequest[CreateRequest](r)
	if err != nil {
		h.logger.Error("Decode create user request", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	user, err := h.userService.Create(r.Context(), model.User{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		h.logger.Error("Create user", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Create user error")

		return
	}

	api.EncodeSuccess(w, http.StatusCreated, CreateResponse{
		ID: user.ID,
	})
}
