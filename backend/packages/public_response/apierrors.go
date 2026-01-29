package public_response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// --- Standard Application Errors ---
// These errors can be returned from business logic layers.
var (
	ErrNotFound       = errors.New("resource not found")
	ErrValidation     = errors.New("validation failed")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrDuplicateEntry = errors.New("duplicate entry")
)

// errorMap maps our standard Go errors to the user-facing ErrorResponse.
var errorMap = map[error]ErrorResponse{
	ErrNotFound:       {Code: "not_found", Description: "The requested resource could not be found."},
	ErrValidation:     {Code: "validation_failed", Description: "The request data is invalid."},
	ErrUnauthorized:   {Code: "unauthorized", Description: "Authentication is required and has failed or has not yet been provided."},
	ErrForbidden:      {Code: "forbidden", Description: "You do not have permission to perform this action."},
	ErrDuplicateEntry: {Code: "duplicate_entry", Description: "The resource you are trying to create already exists."},
}

// statusCodeMap maps our standard Go errors to HTTP status codes.
var statusCodeMap = map[error]int{
	ErrNotFound:       http.StatusNotFound,
	ErrValidation:     http.StatusBadRequest,
	ErrUnauthorized:   http.StatusUnauthorized,
	ErrForbidden:      http.StatusForbidden,
	ErrDuplicateEntry: http.StatusConflict, // 409 Conflict is a good choice for duplicates
}

type ErrorPublicResponse struct {
	Error ErrorResponse `json:"error"`
}

// ErrorResponse is the standard format for API error responses.
type ErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ToError inspects a Go error and writes the appropriate API error response.
// It uses the predefined error maps to find the correct HTTP status and response body.
// If the error is not in the map, it defaults to a 500 Internal Server Error.
func ToError(w http.ResponseWriter, err error) {
	// We iterate because the incoming error might be wrapped (e.g., using fmt.Errorf).
	for key, apiErr := range errorMap {
		if errors.Is(err, key) {
			statusCode := statusCodeMap[key]
			ToErrorResponse(w, statusCode, apiErr.Code, apiErr.Description)
			return
		}
	}

	// If the error is not found in our map, it's an unexpected error.
	ToServerError(w, err)
}

// ToErrorResponse writes a standard client-facing error response.
func ToErrorResponse(w http.ResponseWriter, statusCode int, code, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(
		ErrorPublicResponse{ErrorResponse{
			Code:        code,
			Description: description,
		}})
}

// ToServerError writes a generic 5xx server error response.
// It logs the actual error but returns a generic message to the client.
func ToServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v", err)
	ToErrorResponse(w, http.StatusInternalServerError, "internal_server_error", "An unexpected error occurred.")
}
