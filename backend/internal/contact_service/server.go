package contact_service

import (
	"CRM/internal/contact_service/contracts"
	"CRM/internal/contact_service/service"
	"CRM/packages/public_response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ContactHandlerServer is the HTTP layer.
type ContactHandlerServer struct {
	service *service.ContactService
}

// NewContactHandlerServer creates a new handler and registers its routes.
func NewContactHandlerServer(mux *chi.Mux, svc *service.ContactService) *ContactHandlerServer {
	s := &ContactHandlerServer{
		service: svc,
	}
	RegisterRoutes(mux, s)
	return s
}

// CreateContact handles the HTTP request to create a new contact.
func (h *ContactHandlerServer) CreateContact(w http.ResponseWriter, r *http.Request) {
	// Extract the context from the request.
	ctx := r.Context()

	var req contracts.CreateContactRequest
	// 1. Decode the request.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}

	// 2. Validate the request.
	if err := service.ValidateCreateContactRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}

	// 3. Delegate to the service layer, passing the context.
	response, err := h.service.CreateContact(ctx, req)
	if err != nil {
		// 4. Let the error handler figure out the correct HTTP response.
		public_response.ToError(w, err)
		return
	}

	// 5. Send the successful response.
	public_response.Created(w, response)
}
