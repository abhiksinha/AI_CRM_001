package public_response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON sends a structured JSON response.
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			log.Printf("Failed to encode JSON response: %v", err)
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
