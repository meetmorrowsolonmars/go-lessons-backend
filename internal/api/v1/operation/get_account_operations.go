package operation

import (
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
)

type GetAccountOperationsResponse struct {
	Operations []Operation `json:"operations"`
}

func (h *Handler) GetAccountOperations(w http.ResponseWriter, r *http.Request) {
	const (
		accountIDKey  = "account_id"
		limitKey      = "limit"
		offsetKey     = "offset"
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := api.QueryValueInt64(r, false, limitKey)
	if err != nil {
		h.logger.Error("Get account operations error", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	offset, err := api.QueryValueInt64(r, false, offsetKey)
	if err != nil {
		h.logger.Error("Get account operations error", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	if limit == 0 {
		limit = defaultLimit
	}

	if offset == 0 {
		offset = defaultOffset
	}

	accountID, err := api.PathValueInt64(r, true, accountIDKey)
	if err != nil {
		h.logger.Error("Get account operations error", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	operations, err := h.operationService.GetByAccountID(r.Context(), accountID, limit, offset)
	if err != nil {
		h.logger.Error("Get account operations error", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Get account operations error")

		return
	}

	response := GetAccountOperationsResponse{
		Operations: make([]Operation, 0, len(operations)),
	}

	for _, operation := range operations {
		response.Operations = append(response.Operations, Operation{
			ID:          operation.ID,
			UserID:      operation.UserID,
			AccountID:   operation.AccountID,
			Type:        string(operation.Type),
			CategoryID:  operation.CategoryID,
			Amount:      operation.Amount,
			Description: operation.Description,
			CreateTime:  operation.CreateTime,
		})
	}

	api.EncodeSuccess(w, http.StatusOK, response)
}
