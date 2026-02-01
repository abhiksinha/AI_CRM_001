package contact_service

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes creates a new router group for the v1 API and mounts the contact-related routes.
func RegisterRoutes(router *chi.Mux, handler *ContactHandlerServer) {
	v1 := chi.NewRouter()

	// Here you can add middleware that applies only to this v1 group.
	// For example: v1.Use(authentication.Verifier())

	// Mount the routes for the 'contacts' resource.
	v1.Route("/contacts", func(r chi.Router) {
		r.Post("/", handler.CreateContact)
		r.Get("/", handler.ListContacts) // List all contacts for the owner

		// Routes for a specific contact
		r.Route("/{contactID}", func(r chi.Router) {
			r.Get("/", handler.GetContact)
			r.Put("/", handler.UpdateContact)
			r.Delete("/", handler.DeleteContact)

			// Routes for notes related to a specific contact
			r.Post("/notes", handler.AddNote)
			r.Get("/notes", handler.ListNotes)
		})
	})

	// Mount the v1 router on the main router under the /api/v1 path.
	router.Mount("/api/v1", v1)
}
