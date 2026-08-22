package response

import (
	"encoding/json"
	"net/http"
	"time"
)

type Envelope struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      any         `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

func Write(w http.ResponseWriter, status int, code int, message string, data any, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Code: code, Message: message, Data: data, Timestamp: time.Now().UTC(), RequestID: requestID})
}

func OK(w http.ResponseWriter, data any, requestID string) { Write(w, http.StatusOK, 0, "success", data, requestID) }
func Created(w http.ResponseWriter, data any, requestID string) { Write(w, http.StatusCreated, 0, "success", data, requestID) }
func Error(w http.ResponseWriter, status int, code int, message string, requestID string) { Write(w, status, code, message, nil, requestID) }
