package user

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type GetByIDResponse struct {
	User
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	const pathKey = "user_id"

	userID, err := api.PathValueInt64(r, true, pathKey)
	if err != nil {
		h.logger.Error("Get user by id", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return

	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if errors.Is(err, model.ErrNotFound) {
		h.logger.Info("User not found", slog.Int64("user_id", userID))

		api.EncodeErrorf(w, http.StatusNotFound, "User not found")

		return
	}
	if err != nil {
		h.logger.Error("Get user by id", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Get user by id error")

		return
	}

	api.EncodeSuccess(w, http.StatusOK, GetByIDResponse{
		User: User{
			ID:         user.ID,
			Email:      user.Email,
			FullName:   user.FullName,
			CreateTime: user.CreateTime,
		},
	})
}
