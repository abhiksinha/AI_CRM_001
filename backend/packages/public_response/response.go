package public_response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON sends a structured JSON response.
// It takes a status code and a payload, which can be any serializable object.
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	// Set the content type header.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code.
	w.WriteHeader(statusCode)

	// Encode the payload.
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			// If encoding fails, it's a server-side issue.
			log.Printf("Failed to encode JSON response: %v", err)
			// We can't send another response here as the header is already written.
		}
	}
}

// OK sends a standard 200 OK response with a payload.
func OK(w http.ResponseWriter, payload interface{}) {
	JSON(w, http.StatusOK, payload)
}

// Created sends a standard 201 Created response with a payload.
func Created(w http.ResponseWriter, payload interface{}) {
	JSON(w, http.StatusCreated, payload)
}
