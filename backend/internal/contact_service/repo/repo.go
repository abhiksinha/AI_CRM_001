package repo

import (
	"CRM/internal/contact_service/contracts"
	"CRM/internal/contact_service/model"
	"log"

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

// User is a minimal struct to fetch user data for the response.
type User struct {
	ID        string `gorm:"type:varchar(36);primary_key"`
	FirstName string
	LastName  string
}

// TableName explicitly sets the table name for the User model.
func (u *User) TableName() string {
	return "users"
}

// CreateContact inserts a new contact into the database.
func (r *Repository) CreateContact(contact *model.Contact) error {
	result := r.db.Create(contact)
	return result.Error
}

// ToContactsDbModel converts an API request contract into a database model.
func ToContactsDbModel(req contracts.CreateContactRequest) (*model.Contact, error) {
	return &model.Contact{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		OwnerID:   req.OwnerID,
	}, nil
}

// ToContactsApiResponse converts a database model into an API response,
// enriching it with the owner's name from the database.
func (r *Repository) ToContactsApiResponse(contact *model.Contact) *contracts.CreateContactResponse {
	var owner User
	ownerName := "Unknown" // Default value if owner is not found

	// Find the owner in the users table to get their name.
	if err := r.db.First(&owner, "id = ?", contact.OwnerID).Error; err == nil {
		ownerName = owner.FirstName + " " + owner.LastName
	} else {
		// Log the error if the owner wasn't found, but don't fail the request.
		log.Printf("Could not find owner with ID %s: %v", contact.OwnerID, err)
	}

	return &contracts.CreateContactResponse{
		ID:        contact.ID,
		Name:      contact.FirstName + " " + contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		OwnerID:   contact.OwnerID,
		OwnerName: ownerName, // Use the real name
		CreatedAt: contact.CreatedAt,
	}
}
