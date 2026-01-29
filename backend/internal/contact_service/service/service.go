package service

import (
	"CRM/internal/contact_service/contracts"
	"CRM/internal/contact_service/repo"
	"CRM/packages/logger"
	"CRM/packages/public_response"
	"context"
	"go.uber.org/zap"
)

// ContactService encapsulates the business logic for the contact service.
type ContactService struct {
	repo   *repo.Repository
	logger *zap.Logger
}

// NewContactService creates a new ContactService by applying functional options.
func NewContactService(opts ...Option) *ContactService {
	// Create a new service with default values.
	s := &ContactService{}

	// Apply all the options.
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// CreateContact contains the core logic for creating a new contact.
func (s *ContactService) CreateContact(ctx context.Context, req contracts.CreateContactRequest) (*contracts.CreateContactResponse, error) {
	// Create a logger with the request ID from the context.
	log := logger.FromContext(ctx, s.logger)

	// Convert request to database model.
	newContact, err := repo.ToContactsDbModel(req)
	if err != nil {
		log.Error("Error converting to DB model", zap.Error(err))
		return nil, public_response.ErrValidation
	}

	// Insert into the database.
	if err := s.repo.CreateContact(newContact); err != nil {
		log.Error("Error creating contact in DB", zap.Error(err))
		return nil, err
	}

	// Convert database model to API response.
	response := repo.ToContactsApiResponse(newContact)

	log.Info("Successfully created contact", zap.String("contact_id", newContact.ID))
	return response, nil
}
