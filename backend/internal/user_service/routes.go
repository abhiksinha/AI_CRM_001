package user_service

import "github.com/go-chi/chi/v5"

func RegisterRoutes(router chi.Router, handler *UserHandlerServer) {
	router.Route("/users", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Get("/", handler.ListUsers)
		r.Post("/verify-password", handler.VerifyPassword)

		// API Key Management
		r.Post("/api-keys", handler.CreateApiKey)
		r.Post("/api-keys/match", handler.MatchApiKey)
		r.Delete("/api-keys", handler.ExpireApiKey)

		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", handler.GetUser)
			r.Put("/", handler.UpdateUser)
			r.Delete("/", handler.DeleteUser)
		})
	})
}
