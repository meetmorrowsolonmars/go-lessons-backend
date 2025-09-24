package operation

import (
	"log/slog"
	"net/http"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const operationIDKey = "operation_id"

	operationID, err := api.PathValueUUID(r, true, operationIDKey)
	if err != nil {
		h.logger.Error("Delete operation error", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request body: %s", err)

		return
	}

	err = h.operationService.Delete(r.Context(), operationID)
	// TODO: Handle error.
	if err != nil {
		h.logger.Error("Delete operation", slog.String("error", err.Error()))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Delete operation error")

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
