package server

import (
	"encoding/json"
	"net/http"
)

// APIError is a structured error response.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// writeJSON serializes v as JSON and writes it to w with the given
// status code and appropriate headers.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a structured error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, APIError{
		Code:    status,
		Message: msg,
	})
}

// writeSuccess writes a success response wrapping data in a standard
// envelope: {"ok": true, "data": ...}.
func writeSuccess(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"ok":   true,
		"data": data,
	})
}

// setCORS sets standard CORS headers for API responses.
func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
