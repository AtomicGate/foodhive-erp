package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// helpers for requests
func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func DecodeRequestBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return false
	}
	return true
}

func ParsePaginationParams(r *http.Request) (int, int, error) {
	limit, offset := 10, 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		} else {
			return 0, 0, errors.New("invalid limit parameter")
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		} else {
			return 0, 0, errors.New("invalid offset parameter")
		}
	}

	return limit, offset, nil
}
