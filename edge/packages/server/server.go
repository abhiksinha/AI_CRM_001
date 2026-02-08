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
}

// New creates and configures a new server instance.
func New() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)

	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return s
}

// Router returns the underlying chi router.
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
