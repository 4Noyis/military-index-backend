package utils

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorData `json:"error"`
}

// ErrorData contains error details
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// SendJSON sends a generic JSON response
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendSuccess sends a successful response
func SendSuccess(w http.ResponseWriter, statusCode int, data interface{}, message ...string) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	SendJSON(w, statusCode, response)
}

// SendError sends an error response
func SendError(w http.ResponseWriter, statusCode int, code string, message string, details ...string) {
	response := ErrorResponse{
		Success: false,
		Error: ErrorData{
			Code:    code,
			Message: message,
		},
	}

	if len(details) > 0 {
		response.Error.Details = details[0]
	}

	SendJSON(w, statusCode, response)
}

// SendPaginated sends a paginated response
func SendPaginated(w http.ResponseWriter, data interface{}, page, limit int, total int64) {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	response := PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	SendJSON(w, http.StatusOK, response)
}

// SendCreated sends a 201 Created response
func SendCreated(w http.ResponseWriter, data interface{}, message ...string) {
	SendSuccess(w, http.StatusCreated, data, message...)
}

// SendNoContent sends a 204 No Content response
func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// SendBadRequest sends a 400 Bad Request response
func SendBadRequest(w http.ResponseWriter, message string, details ...string) {
	SendError(w, http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// SendNotFound sends a 404 Not Found response
func SendNotFound(w http.ResponseWriter, message string) {
	SendError(w, http.StatusNotFound, "NOT_FOUND", message)
}

// SendInternalError sends a 500 Internal Server Error response
func SendInternalError(w http.ResponseWriter, message string, details ...string) {
	SendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, details...)
}

// SendUnauthorized sends a 401 Unauthorized response
func SendUnauthorized(w http.ResponseWriter, message string) {
	SendError(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// SendForbidden sends a 403 Forbidden response
func SendForbidden(w http.ResponseWriter, message string) {
	SendError(w, http.StatusForbidden, "FORBIDDEN", message)
}
