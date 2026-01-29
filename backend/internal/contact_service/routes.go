package contact_service

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes creates a new router group for the v1 API and mounts the contact-related routes.
// It takes a ContactHandler as a dependency to map routes to the actual handler methods.
func RegisterRoutes(router *chi.Mux, handler *ContactHandlerServer) {
	// Create a new router group for our /api/v1 endpoints.
	v1 := chi.NewRouter()

	// Here you can add middleware that applies only to this v1 group.
	// For example: v1.Use(authentication.Verifier())

	// Mount the routes for the 'contacts' resource.
	v1.Route("/contacts", func(r chi.Router) {
		r.Post("/", handler.CreateContact)
		// r.Get("/", handler.ListContacts)
		// r.Get("/{contactID}", handler.GetContact)
		// r.Put("/{contactID}", handler.UpdateContact)
		// r.Delete("/{contactID}", handler.DeleteContact)
	})

	// Mount the v1 router on the main router under the /api/v1 path.
	router.Mount("/api/v1", v1)
}
