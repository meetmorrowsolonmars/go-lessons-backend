package account

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type CreateRequest struct {
	Title string `json:"title"`
}

type CreateResponse struct {
	ID int64 `json:"id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := model.AuthClaimsValue(r.Context())

	req, err := api.DecodeRequest[CreateRequest](r)
	if err != nil {
		h.logger.Error("Decode create account request", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	account, err := h.accountService.Create(r.Context(), model.Account{
		UserID: claims.UserID,
		Title:  req.Title,
	})
	if errors.Is(err, model.ErrNotFound) {
		h.logger.Error("Not found",
			slog.Int64("user_id", claims.UserID),
			slog.String("error", err.Error()),
		)

		api.EncodeErrorf(w, http.StatusNotFound, "Not found")

		return
	}
	if err != nil {
		h.logger.Error("Create account", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Create account error")

		return
	}

	api.EncodeSuccess(w, http.StatusCreated, CreateResponse{
		ID: account.ID,
	})
}
