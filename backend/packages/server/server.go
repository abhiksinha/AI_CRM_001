package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server holds the dependencies for our HTTP server.
type Server struct {
	router *chi.Mux
	// We can add other dependencies here later, like a database connection pool.
}

// New creates and configures a new server instance.
func New() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}

	// --- Standard Middleware ---
	s.router.Use(middleware.Logger)    // Log API requests
	s.router.Use(middleware.Recoverer) // Recover from panics
	s.router.Use(middleware.RequestID) // Add a request ID to each request
	s.router.Use(middleware.RealIP)    // Get the real IP from headers

	// --- Health Check Route ---
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return s
}

// Router returns the underlying chi router. This is useful for mounting sub-routers.
func (s *Server) Router() *chi.Mux {
	return s.router
}

// Start runs the HTTP server on a given port.
func (s *Server) Start(port string) {
	fmt.Printf("HTTP server listening on port %s\n", port)
	if err := http.ListenAndServe(port, s.router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
