package operation

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type CreateRequest struct {
	AccountID   int64           `json:"account_id"`
	CategoryID  int64           `json:"category_id,omitempty"`
	Type        string          `json:"type"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}

type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := model.AuthClaimsValue(r.Context())

	req, err := api.DecodeRequest[CreateRequest](r)
	if err != nil {
		h.logger.Error("Decode create operation request", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	// TODO: Validate operation type.
	operation, err := h.operationService.Create(r.Context(), model.Operation{
		UserID:      claims.UserID,
		AccountID:   req.AccountID,
		Type:        model.OperationType(req.Type),
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
	})
	if errors.Is(err, model.ErrNotFound) {
		h.logger.Error("Not found",
			slog.Int64("user_id", claims.UserID),
			slog.Int64("account_id", req.AccountID),
			slog.Int64("category_id", req.CategoryID),
			slog.String("error", err.Error()),
		)

		api.EncodeErrorf(w, http.StatusNotFound, "Not found")

		return
	}
	if err != nil {
		h.logger.Error("Create operation", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Create operation error")

		return
	}

	api.EncodeSuccess(w, http.StatusCreated, CreateResponse{
		ID: operation.ID,
	})
}
