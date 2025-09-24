package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func DecodeRequest[T any](r *http.Request) (T, error) {
	var req T
	return req, json.NewDecoder(r.Body).Decode(&req)
}

func PathValueInt64(r *http.Request, required bool, key string) (int64, error) {
	v := r.PathValue(key)
	if v == "" {
		if required {
			return 0, fmt.Errorf("%s is required", key)
		}

		return 0, nil
	}

	number, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func PathValueUUID(r *http.Request, required bool, key string) (uuid.UUID, error) {
	v := r.PathValue(key)
	if v == "" {
		if required {
			return uuid.Nil, fmt.Errorf("%s is required", key)
		}

		return uuid.Nil, nil
	}

	value, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil, err
	}

	return value, nil
}

func QueryValueInt64(r *http.Request, required bool, key string) (int64, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		if required {
			return 0, fmt.Errorf("%s is required", key)
		}

		return 0, nil
	}

	number, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}
