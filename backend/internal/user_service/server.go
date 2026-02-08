package user_service

import (
	"CRM/internal/user_service/contracts"
	"CRM/internal/user_service/service"
	"CRM/packages/public_response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type UserHandlerServer struct {
	service *service.UserService
}

func NewUserHandlerServer(router chi.Router, svc *service.UserService) *UserHandlerServer {
	s := &UserHandlerServer{service: svc}
	RegisterRoutes(router, s)
	return s
}

// --- API Key Handlers ---

func (h *UserHandlerServer) CreateApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateCreateApiKeyRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	response, err := h.service.CreateApiKey(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.Created(w, response)
}

func (h *UserHandlerServer) MatchApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.MatchApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateMatchApiKeyRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	apiKey, err := h.service.MatchApiKey(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, contracts.ApiKeyMatchResponse{
		IsValid: true,
		UserID:  apiKey.UserID,
	})
}

func (h *UserHandlerServer) ExpireApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.ExpireApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateExpireApiKeyRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	if err := h.service.ExpireApiKey(ctx, req); err != nil {
		public_response.ToError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- User Handlers ---
// (Existing handlers remain unchanged)
func (h *UserHandlerServer) VerifyPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.VerifyPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateVerifyPasswordRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	isValid, err := h.service.VerifyPassword(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, map[string]bool{"is_valid": isValid})
}
func (h *UserHandlerServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateCreateUserRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	response, err := h.service.CreateUser(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.Created(w, response)
}
func (h *UserHandlerServer) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.ListUsersRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "invalid_query_params", err.Error())
		return
	}
	if err := service.ValidateListUsersRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	response, err := h.service.ListUsers(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}
func (h *UserHandlerServer) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userID")
	response, err := h.service.GetUser(ctx, userID)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}
func (h *UserHandlerServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userID")
	var req contracts.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateUpdateUserRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	response, err := h.service.UpdateUser(ctx, userID, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}
func (h *UserHandlerServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userID")
	if err := h.service.DeleteUser(ctx, userID); err != nil {
		public_response.ToError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
