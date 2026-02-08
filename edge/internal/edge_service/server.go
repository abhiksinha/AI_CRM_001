package edge_service

import (
	"bytes"
	"edge/internal/edge_service/contracts"
	"edge/internal/edge_service/service"
	"edge/packages/httpRequest"
	"edge/packages/public_response"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type EdgeHandlerServer struct {
	authService     *service.AuthService
	proxyService    *service.ProxyService
	routeConfig     map[string]RouteConfig
	backendServices map[string]string
}

func NewEdgeHandlerServer(mux *chi.Mux, authSvc *service.AuthService, proxySvc *service.ProxyService, routeConfig map[string]RouteConfig, backendServices map[string]string) *EdgeHandlerServer {
	s := &EdgeHandlerServer{authService: authSvc, proxyService: proxySvc, routeConfig: routeConfig, backendServices: backendServices}
	RegisterRoutes(mux, s)
	return s
}

func (h *EdgeHandlerServer) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}

	response, err := h.authService.Login(ctx, req)
	if err != nil {
		if httpErr, ok := err.(*httpRequest.HTTPError); ok {
			if len(httpErr.Body) == 0 {
				public_response.ToErrorResponse(w, httpErr.StatusCode, "upstream_error", http.StatusText(httpErr.StatusCode))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpErr.StatusCode)
			_, _ = io.Copy(w, bytes.NewReader(httpErr.Body))
			return
		}
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *EdgeHandlerServer) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}

	if err := h.authService.Logout(ctx, req); err != nil {
		public_response.ToError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EdgeHandlerServer) ListSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.SessionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}

	response, err := h.authService.ListSessions(ctx, req)
	if err != nil {
		if httpErr, ok := err.(*httpRequest.HTTPError); ok {
			if len(httpErr.Body) == 0 {
				public_response.ToErrorResponse(w, httpErr.StatusCode, "upstream_error", http.StatusText(httpErr.StatusCode))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpErr.StatusCode)
			_, _ = io.Copy(w, bytes.NewReader(httpErr.Body))
			return
		}
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *EdgeHandlerServer) ExpireSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.ExpireSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}

	if err := h.authService.ExpireSession(ctx, req); err != nil {
		public_response.ToError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EdgeHandlerServer) GetToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.GetTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}

	response, err := h.authService.GetToken(ctx, req)
	if err != nil {
		if httpErr, ok := err.(*httpRequest.HTTPError); ok {
			if len(httpErr.Body) == 0 {
				public_response.ToErrorResponse(w, httpErr.StatusCode, "upstream_error", http.StatusText(httpErr.StatusCode))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpErr.StatusCode)
			_, _ = io.Copy(w, bytes.NewReader(httpErr.Body))
			return
		}
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}

func (h *EdgeHandlerServer) Proxy(w http.ResponseWriter, r *http.Request) {
	routePattern := chi.RouteContext(r.Context()).RoutePattern()
	key := r.Method + " " + routePattern
	cfg, ok := h.routeConfig[key]
	if !ok && routePattern != "" {
		key = r.Method + " /api/v1" + routePattern
		cfg, ok = h.routeConfig[key]
	}
	if !ok && routePattern != "" {
		key = r.Method + " " + strings.TrimPrefix(routePattern, "/api/v1")
		cfg, ok = h.routeConfig[key]
	}
	if !ok {
		public_response.ToError(w, public_response.ErrNotFound)
		return
	}

	authRequired := cfg.AuthType != AuthNone

	var userID string
	if authRequired {
		token, err := parseBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			public_response.ToError(w, public_response.ErrUnauthorized)
			return
		}
		userID, err = h.proxyService.VerifyToken(r.Context(), token)
		if err != nil {
			public_response.ToError(w, public_response.ErrUnauthorized)
			return
		}
	}

	baseURL := h.backendServices[cfg.Service]
	if baseURL == "" {
		public_response.ToServerError(w, fmt.Errorf("missing backend base url for service: %s", cfg.Service))
		return
	}

	if err := h.proxyService.Forward(w, r, baseURL, cfg.BackendPath, userID); err != nil {
		public_response.ToServerError(w, err)
		return
	}
}

func parseBearerToken(header string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", fmt.Errorf("missing bearer token")
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", fmt.Errorf("missing bearer token")
	}
	return token, nil
}
