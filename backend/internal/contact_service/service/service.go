package service

import (
	"CRM/internal/contact_service/contracts"
	"CRM/internal/contact_service/model"
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
	repo      *repo.Repository
	notesRepo *repo.NotesRepository
	logger    *zap.Logger
}

// NewContactService creates a new ContactService.
func NewContactService(repo *repo.Repository, notesRepo *repo.NotesRepository, logger *zap.Logger) *ContactService {
	return &ContactService{
		repo:      repo,
		notesRepo: notesRepo,
		logger:    logger,
	}
}

func (s *ContactService) CreateContact(ctx context.Context, req contracts.CreateContactRequest) (*contracts.ContactResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	existingContact, err := s.repo.GetByQuery(map[string]interface{}{"email": req.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("Error checking for existing contact", zap.Error(err))
		return nil, err
	}
	if existingContact != nil {
		log.Warn("Attempted to create a duplicate contact", zap.String("email", req.Email))
		// Return the specific error struct with the conflicting user ID.
		return nil, &public_response.DuplicateEntryError{UserID: existingContact.OwnerID}
	}
	newContact, err := repo.ToContactsDbModel(req)
	if err != nil {
		log.Error("Error converting to DB model", zap.Error(err))
		return nil, public_response.ErrValidation
	}
	err = s.repo.ExecTxn(ctx, func(txnRepo *repo.Repository) error {
		return txnRepo.CreateContact(newContact)
	})
	if err != nil {
		log.Error("Error during contact creation transaction", zap.Error(err))
		return nil, err
	}
	response := s.repo.ToContactsApiResponse(newContact)
	log.Info("Successfully created contact", zap.String("contact_id", newContact.ID))
	return response, nil
}

func (s *ContactService) GetContactByID(ctx context.Context, id, ownerID string) (*contracts.ContactResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	contact, err := s.repo.GetByIDAndOwner(id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Contact not found", zap.String("id", id), zap.String("owner_id", ownerID))
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching contact", zap.Error(err))
		return nil, err
	}
	return s.repo.ToContactsApiResponse(contact), nil
}

func (s *ContactService) ListContactsByOwner(ctx context.Context, ownerID string, req contracts.ListContactsRequest) (*contracts.ListContactsResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	contacts, total, err := s.repo.ListByOwner(ownerID, req)
	if err != nil {
		log.Error("Error listing contacts", zap.Error(err))
		return nil, err
	}
	response := make([]contracts.ContactResponse, len(contacts))
	for i, c := range contacts {
		response[i] = *s.repo.ToContactsApiResponse(&c)
	}
	return &contracts.ListContactsResponse{Data: response, TotalCount: total}, nil
}

func (s *ContactService) UpdateContact(ctx context.Context, id, ownerID string, req contracts.UpdateContactRequest) (*contracts.ContactResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	contact, err := s.repo.GetByIDAndOwner(id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Contact not found for update", zap.String("id", id), zap.String("owner_id", ownerID))
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching contact for update", zap.Error(err))
		return nil, err
	}

	// Apply updates from the request.
	if req.FirstName != nil {
		contact.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		contact.LastName = *req.LastName
	}
	if req.Email != nil {
		contact.Email = *req.Email
	}
	if req.Phone != nil {
		contact.Phone = *req.Phone
	}
	if req.OwnerID != nil {
		contact.OwnerID = *req.OwnerID
	}

	err = s.repo.ExecTxn(ctx, func(txnRepo *repo.Repository) error {
		return txnRepo.UpdateContact(contact)
	})
	if err != nil {
		log.Error("Error during contact update transaction", zap.Error(err))
		return nil, err
	}

	return s.repo.ToContactsApiResponse(contact), nil
}

func (s *ContactService) DeleteContact(ctx context.Context, id, ownerID string) error {
	log := logger.FromContext(ctx, s.logger)
	err := s.repo.ExecTxn(ctx, func(txnRepo *repo.Repository) error {
		rowsAffected, err := txnRepo.DeleteByIDAndOwner(id, ownerID)
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return public_response.ErrNotFound
		}
		return nil
	})
	if err != nil {
		log.Error("Error during contact deletion transaction", zap.Error(err))
		return err
	}
	log.Info("Successfully deleted contact", zap.String("id", id), zap.String("owner_id", ownerID))
	return nil
}

func (s *ContactService) AddNoteToContact(ctx context.Context, contactID, ownerID, authorID string, req contracts.AddNoteRequest) (*contracts.NoteResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	if _, err := s.repo.GetByIDAndOwner(contactID, ownerID); err != nil {
		log.Warn("Attempted to add note to non-existent or unauthorized contact", zap.String("contact_id", contactID), zap.String("owner_id", ownerID))
		return nil, public_response.ErrNotFound
	}

	note := &model.Note{
		ContactID: contactID,
		AuthorID:  authorID,
		Content:   req.Content,
	}

	if err := s.notesRepo.CreateNote(note); err != nil {
		log.Error("Error creating note", zap.Error(err))
		return nil, err
	}

	return &contracts.NoteResponse{
		ID:        note.ID,
		Content:   note.Content,
		AuthorID:  note.AuthorID,
		CreatedAt: note.CreatedAt,
	}, nil
}

func (s *ContactService) ListNotesForContact(ctx context.Context, contactID, ownerID string) (*contracts.ListNotesResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	if _, err := s.repo.GetByIDAndOwner(contactID, ownerID); err != nil {
		log.Warn("Attempted to list notes for non-existent or unauthorized contact", zap.String("contact_id", contactID), zap.String("owner_id", ownerID))
		return nil, public_response.ErrNotFound
	}

	notes, err := s.notesRepo.ListByContactID(contactID)
	if err != nil {
		log.Error("Error listing notes", zap.Error(err))
		return nil, err
	}

	response := make([]contracts.NoteResponse, len(notes))
	for i, n := range notes {
		response[i] = contracts.NoteResponse{
			ID:        n.ID,
			Content:   n.Content,
			AuthorID:  n.AuthorID,
			CreatedAt: n.CreatedAt,
		}
	}

	return &contracts.ListNotesResponse{Data: response, TotalCount: int64(len(response))}, nil
}
