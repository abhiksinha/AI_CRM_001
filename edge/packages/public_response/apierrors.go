package public_response

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// --- Standard Application Errors ---
var (
	ErrNotFound       = errors.New("resource not found")
	ErrValidation     = errors.New("validation failed")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrSessionLimit   = errors.New("session limit exceeded")
)

// DuplicateEntryError is a custom error type that includes the ID of the conflicting user.
type DuplicateEntryError struct {
	UserID string
}

func (e *DuplicateEntryError) Error() string {
	return fmt.Sprintf("the customer is already assigned to user id: %s", e.UserID)
}

// errorMap maps our standard Go errors to the user-facing ErrorResponse.
var errorMap = map[error]ErrorResponse{
	ErrNotFound:       {Code: "not_found", Description: "The requested resource could not be found."},
	ErrValidation:     {Code: "validation_failed", Description: "The request data is invalid."},
	ErrUnauthorized:   {Code: "unauthorized", Description: "Authentication is required and has failed or has not yet been provided."},
	ErrForbidden:      {Code: "forbidden", Description: "You do not have permission to perform this action."},
	ErrDuplicateEntry: {Code: "duplicate_entry", Description: "The resource you are trying to create already exists."},
	ErrSessionLimit:   {Code: "session_limit_exceeded", Description: "Maximum active sessions reached."},
}

// statusCodeMap maps our standard Go errors to HTTP status codes.
var statusCodeMap = map[error]int{
	ErrNotFound:       http.StatusNotFound,
	ErrValidation:     http.StatusBadRequest,
	ErrUnauthorized:   http.StatusUnauthorized,
	ErrForbidden:      http.StatusForbidden,
	ErrDuplicateEntry: http.StatusConflict,
	ErrSessionLimit:   http.StatusTooManyRequests,
}

// ErrorResponse is the standard format for API error responses.
type ErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ToError inspects a Go error and writes the appropriate API error response.
func ToError(w http.ResponseWriter, err error) {
	var dupErr *DuplicateEntryError
	if errors.As(err, &dupErr) {
		ToErrorResponse(w, http.StatusConflict, "duplicate_entry", dupErr.Error())
		return
	}

	for key, apiErr := range errorMap {
		if errors.Is(err, key) {
			statusCode := statusCodeMap[key]
			JSON(w, statusCode, apiErr)
			return
		}
	}
	ToServerError(w, err)
}

// ToErrorResponse writes a standard client-facing error response.
func ToErrorResponse(w http.ResponseWriter, statusCode int, code, description string) {
	JSON(w, statusCode, ErrorResponse{Code: code, Description: description})
}

// ToServerError writes a generic 5xx server error response.
func ToServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v", err)
	errorResponse := ErrorResponse{Code: "internal_server_error", Description: "An unexpected error occurred."}
	JSON(w, http.StatusInternalServerError, errorResponse)
}
