package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Success(w http.ResponseWriter, message string, data interface{}) {
	response := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusOK, response)
}

func Error(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    statusCode,
	}
	JSON(w, statusCode, response)
}
