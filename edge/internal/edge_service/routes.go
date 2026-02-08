package edge_service

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type RouteConfig struct {
	BackendPath string
	AuthType    string
	Method      string
	Service     string
}

const (
	AuthNone  = "none"
	AuthToken = "token"
)

func BuildRouteConfig() map[string]RouteConfig {
	return map[string]RouteConfig{
		// GET
		"GET /v1/contacts":                   {BackendPath: "/api/v1/contacts", AuthType: AuthToken, Method: http.MethodGet, Service: "backend_service"},
		"GET /v1/contacts/{contactID}":       {BackendPath: "/api/v1/contacts/{contactID}", AuthType: AuthToken, Method: http.MethodGet, Service: "backend_service"},
		"GET /v1/contacts/{contactID}/notes": {BackendPath: "/api/v1/contacts/{contactID}/notes", AuthType: AuthToken, Method: http.MethodGet, Service: "backend_service"},
		"GET /v1/deals":                      {BackendPath: "/api/v1/deals", AuthType: AuthToken, Method: http.MethodGet, Service: "backend_service"},
		"GET /v1/deals/{dealID}":             {BackendPath: "/api/v1/deals/{dealID}", AuthType: AuthToken, Method: http.MethodGet, Service: "backend_service"},
		"GET /v1/deals/{dealID}/tasks":       {BackendPath: "/api/v1/deals/{dealID}/tasks", AuthType: AuthToken, Method: http.MethodGet, Service: "backend_service"},
		"GET /v1/users":                      {BackendPath: "/api/v1/users", AuthType: AuthToken, Method: http.MethodGet, Service: "backend_service"},
		"GET /v1/users/{userID}":             {BackendPath: "/api/v1/users/{userID}", AuthType: AuthToken, Method: http.MethodGet, Service: "backend_service"},

		// POST
		"POST /v1/contacts":                   {BackendPath: "/api/v1/contacts", AuthType: AuthToken, Method: http.MethodPost, Service: "backend_service"},
		"POST /v1/contacts/{contactID}/notes": {BackendPath: "/api/v1/contacts/{contactID}/notes", AuthType: AuthToken, Method: http.MethodPost, Service: "backend_service"},
		"POST /v1/deals":                      {BackendPath: "/api/v1/deals", AuthType: AuthToken, Method: http.MethodPost, Service: "backend_service"},
		"POST /v1/deals/{dealID}/tasks":       {BackendPath: "/api/v1/deals/{dealID}/tasks", AuthType: AuthToken, Method: http.MethodPost, Service: "backend_service"},
		"POST /v1/users":                      {BackendPath: "/api/v1/users", AuthType: AuthNone, Method: http.MethodPost, Service: "backend_service"},
		"POST /v1/users/verify-password":      {BackendPath: "/api/v1/users/verify-password", AuthType: AuthToken, Method: http.MethodPost, Service: "backend_service"},
		"POST /v1/users/api-keys":             {BackendPath: "/api/v1/users/api-keys", AuthType: AuthToken, Method: http.MethodPost, Service: "backend_service"},
		"POST /v1/users/api-keys/match":       {BackendPath: "/api/v1/users/api-keys/match", AuthType: AuthToken, Method: http.MethodPost, Service: "backend_service"},

		// PUT
		"PUT /v1/contacts/{contactID}": {BackendPath: "/api/v1/contacts/{contactID}", AuthType: AuthToken, Method: http.MethodPut, Service: "backend_service"},
		"PUT /v1/deals/{dealID}":       {BackendPath: "/api/v1/deals/{dealID}", AuthType: AuthToken, Method: http.MethodPut, Service: "backend_service"},
		"PUT /v1/users/{userID}":       {BackendPath: "/api/v1/users/{userID}", AuthType: AuthToken, Method: http.MethodPut, Service: "backend_service"},

		// DELETE
		"DELETE /v1/contacts/{contactID}": {BackendPath: "/api/v1/contacts/{contactID}", AuthType: AuthToken, Method: http.MethodDelete, Service: "backend_service"},
		"DELETE /v1/deals/{dealID}":       {BackendPath: "/api/v1/deals/{dealID}", AuthType: AuthToken, Method: http.MethodDelete, Service: "backend_service"},
		"DELETE /v1/users/api-keys":       {BackendPath: "/api/v1/users/api-keys", AuthType: AuthToken, Method: http.MethodDelete, Service: "backend_service"},
		"DELETE /v1/users/{userID}":       {BackendPath: "/api/v1/users/{userID}", AuthType: AuthToken, Method: http.MethodDelete, Service: "backend_service"},
	}
}

func RegisterRoutes(router *chi.Mux, handler *EdgeHandlerServer) {
	v1 := chi.NewRouter()
	v1.Post("/login", handler.Login)
	v1.Post("/logout", handler.Logout)
	v1.Post("/sessions", handler.ListSessions)
	v1.Post("/sessions/expire", handler.ExpireSession)
	v1.Post("/token", handler.GetToken)
	router.Mount("/v1", v1)

	for key, cfg := range handler.routeConfig {
		path := strings.TrimPrefix(key, cfg.Method+" ")
		if !strings.HasPrefix(path, "/v1/") {
			continue
		}
		router.Method(cfg.Method, path, http.HandlerFunc(handler.Proxy))
	}
}
