package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func RespondJSON(w http.ResponseWriter, status int, payload any) {
	response := APIResponse{
		Data: payload,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error (in a real app, use a proper logger)
	}
}

func RespondError(w http.ResponseWriter, status int, message string) {
	response := APIResponse{
		Error: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error
	}
}

func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
