package handlers

import (
	"encoding/json"
	"net/http"
)

type contextKey string

const ParamsKey contextKey = "params"

func Params(r *http.Request) map[string]string {
	if params, ok := r.Context().Value(ParamsKey).(map[string]string); ok {
		return params
	}
	return nil
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "error",
		"error": map[string]string{
			"message": message,
		},
	})
}
