package deal_service

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router *chi.Mux, handler *DealHandlerServer) {
	v1 := chi.NewRouter()
	v1.Route("/deals", func(r chi.Router) {
		r.Post("/", handler.CreateDeal)
		r.Get("/", handler.ListDeals)

		r.Route("/{dealID}", func(r chi.Router) {
			r.Get("/", handler.GetDeal)
			r.Put("/", handler.UpdateDeal)
			r.Delete("/", handler.DeleteDeal)

			r.Post("/tasks", handler.CreateTask)
			r.Get("/tasks", handler.ListTasks)
		})
	})
	router.Mount("/api/v1", v1)
}
