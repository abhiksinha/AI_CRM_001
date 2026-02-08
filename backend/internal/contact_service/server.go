package contact_service

import (
	"CRM/internal/contact_service/contracts"
	"CRM/internal/contact_service/service"
	"CRM/packages/public_response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// ContactHandlerServer is the HTTP layer.
type ContactHandlerServer struct {
	service *service.ContactService
}

// NewContactHandlerServer creates a new handler and registers its routes.
func NewContactHandlerServer(router chi.Router, svc *service.ContactService) *ContactHandlerServer {
	s := &ContactHandlerServer{
		service: svc,
	}
	RegisterRoutes(router, s)
	return s
}

func getOwnerID(r *http.Request) string {
	return r.Header.Get("X-User-ID")
}

func (h *ContactHandlerServer) CreateContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.CreateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	req.OwnerID = getOwnerID(r)
	if req.OwnerID == "" {
		public_response.ToError(w, public_response.ErrUnauthorized)
		return
	}
	if err := service.ValidateCreateContactRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	response, err := h.service.CreateContact(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.Created(w, response)
}

func (h *ContactHandlerServer) GetContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contactID := chi.URLParam(r, "contactID")
	ownerID := getOwnerID(r)

	response, err := h.service.GetContactByID(ctx, contactID, ownerID)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *ContactHandlerServer) ListContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ownerID := getOwnerID(r)

	var req contracts.ListContactsRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "invalid_query_params", err.Error())
		return
	}

	response, err := h.service.ListContactsByOwner(ctx, ownerID, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *ContactHandlerServer) UpdateContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contactID := chi.URLParam(r, "contactID")
	ownerID := getOwnerID(r)

	var req contracts.UpdateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateUpdateContactRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}

	response, err := h.service.UpdateContact(ctx, contactID, ownerID, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *ContactHandlerServer) DeleteContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contactID := chi.URLParam(r, "contactID")
	ownerID := getOwnerID(r)

	if err := h.service.DeleteContact(ctx, contactID, ownerID); err != nil {
		public_response.ToError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ContactHandlerServer) AddNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contactID := chi.URLParam(r, "contactID")
	ownerID := getOwnerID(r)
	authorID := getOwnerID(r)

	var req contracts.AddNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateAddNoteRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}

	response, err := h.service.AddNoteToContact(ctx, contactID, ownerID, authorID, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.Created(w, response)
}

func (h *ContactHandlerServer) ListNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contactID := chi.URLParam(r, "contactID")
	ownerID := getOwnerID(r)

	response, err := h.service.ListNotesForContact(ctx, contactID, ownerID)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}
