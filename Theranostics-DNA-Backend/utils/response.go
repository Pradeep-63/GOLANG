package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	// Status  int         `json:"status"`
	// Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Will omit if Data is nil
}

func JSONResponse(w http.ResponseWriter, status int, success bool, message string, data interface{}) {
	if data == nil {
		data = struct{}{}
	}

	response := Response{
		// Status:  status,
		// Success: success,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
