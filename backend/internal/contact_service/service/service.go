package service

import (
	"CRM/internal/contact_service/contracts"
	"CRM/internal/contact_service/repo"
	"CRM/packages/logger"
	"CRM/packages/public_response"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
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
	log := logger.FromContext(ctx, s.logger)

	// 1. Check for duplicates before starting a transaction.
	existingContact, err := s.repo.GetByQuery(map[string]interface{}{"email": req.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("Error checking for existing contact", zap.Error(err))
		return nil, err
	}
	if existingContact != nil {
		log.Warn("Attempted to create a duplicate contact", zap.String("email", req.Email))
		return nil, public_response.ErrDuplicateEntry
	}

	// 2. Map the request to the database model.
	newContact, err := repo.ToContactsDbModel(req)
	if err != nil {
		log.Error("Error converting to DB model", zap.Error(err))
		return nil, public_response.ErrValidation
	}

	// Insert into the database.
	err = s.repo.ExecTxn(ctx, func(repo *repo.Repository) error {
		if err := s.repo.CreateContact(newContact); err != nil {
			log.Error("Error creating contact in DB", zap.Error(err))
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. Convert the successful model to an API response.
	response := s.repo.ToContactsApiResponse(newContact)

	log.Info("Successfully created contact", zap.String("contact_id", newContact.ID))
	return response, nil
}
