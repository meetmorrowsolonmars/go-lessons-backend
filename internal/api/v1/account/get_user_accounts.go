package account

import (
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type GetUserAccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

func (h *Handler) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	claims, _ := model.AuthClaimsValue(r.Context())

	accounts, err := h.accountService.GetAccountsByUserID(r.Context(), claims.UserID)
	if err != nil {
		h.logger.Error("Get user accounts", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Get user accounts error")

		return
	}

	response := GetUserAccountsResponse{
		Accounts: make([]Account, 0, len(accounts)),
	}

	for _, account := range accounts {
		response.Accounts = append(response.Accounts, Account{
			ID:         account.ID,
			UserID:     account.UserID,
			Title:      account.Title,
			IsDefault:  account.IsDefault,
			CreateTime: account.CreateTime,
		})
	}

	api.EncodeSuccess(w, http.StatusOK, response)
}
