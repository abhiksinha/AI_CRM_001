package contact_service

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers contact-related routes on the provided router.
func RegisterRoutes(router chi.Router, handler *ContactHandlerServer) {
	// Here you can add middleware that applies only to this group.
	// For example: router.Use(authentication.Verifier())

	router.Route("/contacts", func(r chi.Router) {
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
}
