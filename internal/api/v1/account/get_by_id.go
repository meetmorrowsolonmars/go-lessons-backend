package account

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type GetByIDResponse struct {
	Account
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	const accountIDKey = "account_id"

	accountID, err := api.PathValueInt64(r, true, accountIDKey)
	if err != nil {
		h.logger.Error("Get account by id", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	account, err := h.accountService.GetByID(r.Context(), accountID)
	if errors.Is(err, model.ErrNotFound) {
		h.logger.Error("Account not found", slog.Int64("account_id", accountID))

		api.EncodeErrorf(w, http.StatusNotFound, "Account not found")

		return
	}
	if err != nil {
		h.logger.Error("Get account by id", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Get account error")

		return
	}

	api.EncodeSuccess(w, http.StatusOK, GetByIDResponse{
		Account: Account{
			ID:         account.ID,
			UserID:     account.UserID,
			Title:      account.Title,
			IsDefault:  account.IsDefault,
			CreateTime: account.CreateTime,
		},
	})
}
