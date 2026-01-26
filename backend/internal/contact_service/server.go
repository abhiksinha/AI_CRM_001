package contact_service

import (
	// Import the generated RPC code
	"context"
	"time"

	"CRM/rpc/proto/contacts"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the ContactService gRPC service.
type Server struct {
	// This ensures we implement all methods of the interface,
	// even if we haven't explicitly written them yet.
	rpc.UnimplementedContactServiceServer
}

// NewServer creates a new instance of the contact service server.
func NewServer() *Server {
	return &Server{}
}

// CreateContact is the implementation of the RPC method to create a new contact.
func (s *Server) CreateContact(ctx context.Context, req *rpc.CreateContactRequest) (*rpc.CreateContactResponse, error) {
	// TODO: Add database logic here to insert the new contact.
	// TODO: Add input validation.

	// For now, we'll return a hardcoded response to confirm the wiring is correct.
	newID := uuid.New().String()
	now := timestamppb.New(time.Now())

	// Create a new contact object based on the request.
	contact := &rpc.Contact{
		Id:        newID,
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Email:     req.GetEmail(),
		Phone:     req.GetPhone(),
		OwnerId:   req.GetOwnerId(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create the response.
	res := &rpc.CreateContactResponse{
		Contact: contact,
	}

	return res, nil
}
