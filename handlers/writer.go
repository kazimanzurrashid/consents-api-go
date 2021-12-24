package handlers

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, statusCode int, errors ...string) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)

	res, _ := json.Marshal(struct {
		Errors []string `json:"errors"`
	}{
		Errors: errors,
	})

	_, _ = w.Write(res)
}

func writeSuccess(w http.ResponseWriter, statusCode int, payload interface{}) {
	res, err := json.Marshal(payload)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(res)
}
