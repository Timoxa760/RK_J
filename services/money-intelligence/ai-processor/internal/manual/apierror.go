package manual

import (
	"encoding/json"
	"net/http"
)

// APIError — ошибка с кодом для клиента.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"error"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIError{Code: code, Message: message})
}

func respondError(w http.ResponseWriter, code int, err error) {
	status := http.StatusBadRequest
	if code == 500 {
		status = http.StatusInternalServerError
	}
	if apiErr, ok := err.(*APIError); ok {
		writeAPIError(w, status, apiErr.Code, apiErr.Message)
		return
	}
	writeAPIError(w, status, "error", err.Error())
}
