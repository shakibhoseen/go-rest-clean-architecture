package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"error message description"`
}

type PaginatedResponse struct {
	TotalRecords int `json:"total_records"`
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	Limit        int `json:"limit"`
	Data         any `json:"data"`
}

func JSON(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}

func ErrorJSON(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIResponse{
		Success: false,
		Error:   errMsg,
	}
	json.NewEncoder(w).Encode(resp)
}
