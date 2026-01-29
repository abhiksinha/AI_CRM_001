package repo

import (
	"CRM/internal/contact_service/contracts"
	"CRM/internal/contact_service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository provides access to the database for the contact service.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateContact inserts a new contact into the database.
func (r *Repository) CreateContact(contact *model.Contact) error {
	result := r.db.Create(contact)
	return result.Error
}

// ToContactsDbModel converts an API request contract into a database model.
func ToContactsDbModel(req contracts.CreateContactRequest) (*model.Contact, error) {
	ownerID, err := uuid.Parse(req.OwnerID)
	if err != nil {
		// Return an error if the OwnerID is not a valid UUID.
		return nil, err
	}

	return &model.Contact{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		OwnerID:   ownerID,
	}, nil
}

// ToContactsApiResponse converts a database model into an API response contract.
func ToContactsApiResponse(contact *model.Contact) *contracts.CreateContactResponse {
	// TODO: Get OwnerName from the database via a join or a separate query.
	ownerName := "Dummy Owner Name"

	return &contracts.CreateContactResponse{
		ID:        contact.ID,
		Name:      contact.FirstName + " " + contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		OwnerID:   contact.OwnerID.String(),
		OwnerName: ownerName,
		CreatedAt: contact.CreatedAt,
	}
}
