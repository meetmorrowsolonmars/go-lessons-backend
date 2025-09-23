package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func EncodeErrorf(w http.ResponseWriter, code int, format string, a ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(defaultErrorResponse{
		Message: fmt.Sprintf(format, a...),
	})
}

type defaultErrorResponse struct {
	Message string `json:"message"`
}

func EncodeSuccess[T any](w http.ResponseWriter, code int, response T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(response)
}
