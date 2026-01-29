package main

import (
	"CRM/internal/contact_service"
	"CRM/packages/server"
)

func main() {
	// Create a new server instance from our server package.
	srv := server.New()

	// --- Register API Routes ---
	// 1. Create the handler for the contact service.
	contact_service.NewContactHandlerServer(srv.Router())
	// 2. Register the v1 API routes, passing the handler to the router function.

	// Start the server on port 8080.
	srv.Start(":8080")
}
