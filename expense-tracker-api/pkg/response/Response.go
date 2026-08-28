package response

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, r HTTPResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(r)
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, HTTPResponse{
		Success: false,
		Message: message,
		Error:   message,
	})
}