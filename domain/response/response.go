package response

import (
	"encoding/json"
	"net/http"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(message string, data interface{}) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string, err error) *Response {
	response := &Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return response
}

// WriteJSON writes a response as JSON to the HTTP response writer
func (r *Response) WriteJSON(w http.ResponseWriter, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(r)
}

// WriteSuccess writes a success response
func WriteSuccess(w http.ResponseWriter, message string, data interface{}) error {
	return SuccessResponse(message, data).WriteJSON(w, http.StatusOK)
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, statusCode int, message string, err error) error {
	return ErrorResponse(message, err).WriteJSON(w, statusCode)
}

// WriteBadRequest writes a bad request error response
func WriteBadRequest(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusBadRequest, message, nil)
}

// WriteInternalError writes an internal server error response
func WriteInternalError(w http.ResponseWriter, message string, err error) error {
	return WriteError(w, http.StatusInternalServerError, message, err)
}

// WriteNotFound writes a not found error response
func WriteNotFound(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusNotFound, message, nil)
}
