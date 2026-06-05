package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message,omitempty"`
	Code      int    `json:"code"`
	RequestID string `json:"request_id,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, errMsg string) error {
	reqID := w.Header().Get("X-Request-ID")
	return JSON(w, status, ErrorResponse{
		Error:     errMsg,
		Code:      status,
		RequestID: reqID,
	})
}

func Success(w http.ResponseWriter, message string, data interface{}) error {
	return JSON(w, http.StatusOK, SuccessResponse{Message: message, Data: data})
}
