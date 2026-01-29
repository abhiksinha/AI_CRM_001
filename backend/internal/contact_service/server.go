package contact_service

import (
	"CRM/internal/contact_service/contracts" // Corrected import path
	"CRM/packages/public_response"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ContactHandlerServer holds the dependencies for the contact handlers, like a database connection.
type ContactHandlerServer struct {
	// db *sql.DB
}

// NewContactHandlerServer creates a new ContactHandler and registers its routes.
func NewContactHandlerServer(mux *chi.Mux) *ContactHandlerServer {
	s := &ContactHandlerServer{}
	RegisterRoutes(mux, s)
	return s
}

// CreateContact handles the HTTP request to create a new contact.
func (h *ContactHandlerServer) CreateContact(w http.ResponseWriter, r *http.Request) {
	var req contracts.CreateContactRequest
	// Decode the JSON request body.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}

	// TODO: Add input validation.
	// TODO: Add database logic to insert the new contact and get owner name.

	// For now, create a dummy response using the new struct.
	response := contracts.CreateContactResponse{
		ID:        uuid.New().String(), // Generate a new UUID for the contact.
		Name:      req.FirstName + " " + req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		OwnerID:   req.OwnerID,
		OwnerName: "Dummy Owner Name", // This would come from a DB join.
		CreatedAt: time.Now(),
	}

	// Use the new helper to send a 201 Created response.
	public_response.Created(w, response)
}
