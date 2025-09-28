package operation

import (
	"log/slog"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
)

type UpdateRequest struct {
	Amount      decimal.Decimal `json:"amount"`
	CategoryID  int64           `json:"category_id,omitempty"`
	Description string          `json:"description"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	const operationIDKey = "operation_id"

	operationID, err := api.PathValueUUID(r, true, operationIDKey)
	if err != nil {
		h.logger.Error("Update operation error", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	req, err := api.DecodeRequest[UpdateRequest](r)
	if err != nil {
		h.logger.Error("Decode update operation request", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	err = h.operationService.Update(r.Context(), operationID, req.Amount, req.CategoryID, req.Description)
	// TODO: Handle error.
	if err != nil {
		h.logger.Error("Update operation", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Update operation error")

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
