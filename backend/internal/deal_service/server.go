package deal_service

import (
	"CRM/internal/deal_service/contracts"
	"CRM/internal/deal_service/service"
	"CRM/packages/public_response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type DealHandlerServer struct {
	service *service.DealService
}

func NewDealHandlerServer(router chi.Router, svc *service.DealService) *DealHandlerServer {
	s := &DealHandlerServer{service: svc}
	RegisterRoutes(router, s)
	return s
}

func getOwnerID(r *http.Request) string {
	return r.Header.Get("X-User-ID")
}

func (h *DealHandlerServer) CreateDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.CreateDealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateCreateDealRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	req.OwnerID = getOwnerID(r)
	response, err := h.service.CreateDeal(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.Created(w, response)
}

func (h *DealHandlerServer) GetDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dealID := chi.URLParam(r, "dealID")
	ownerID := getOwnerID(r)
	response, err := h.service.GetDealByID(ctx, dealID, ownerID)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *DealHandlerServer) ListDeals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ownerID := getOwnerID(r)
	var req contracts.ListDealsRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "invalid_query_params", err.Error())
		return
	}
	response, err := h.service.ListDeals(ctx, ownerID, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *DealHandlerServer) UpdateDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dealID := chi.URLParam(r, "dealID")
	ownerID := getOwnerID(r)
	var req contracts.UpdateDealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateUpdateDealRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	response, err := h.service.UpdateDeal(ctx, dealID, ownerID, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *DealHandlerServer) DeleteDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dealID := chi.URLParam(r, "dealID")
	ownerID := getOwnerID(r)
	if err := h.service.DeleteDeal(ctx, dealID, ownerID); err != nil {
		public_response.ToError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *DealHandlerServer) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dealID := chi.URLParam(r, "dealID")
	ownerID := getOwnerID(r)
	var req contracts.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateCreateTaskRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}
	response, err := h.service.CreateTask(ctx, dealID, ownerID, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.Created(w, response)
}

func (h *DealHandlerServer) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dealID := chi.URLParam(r, "dealID")
	ownerID := getOwnerID(r)
	response, err := h.service.ListTasks(ctx, dealID, ownerID)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}
