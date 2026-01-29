package contact_service

import (
	"CRM/internal/contact_service/contracts"
	"CRM/packages/public_response" // Import the new package
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	// TODO: Add database logic to insert the new contact.

	// Use the new helper to send a 201 Created response.
	public_response.Created(w, map[string]string{"message": "Contact created successfully"})
}
