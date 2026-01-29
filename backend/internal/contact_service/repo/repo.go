package repo

import (
	"CRM/internal/contact_service/contracts"
	"CRM/internal/contact_service/model"
	"context"
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

// ExecTxn executes the given function within a database transaction.
// If the function returns an error, the transaction is rolled back. Otherwise, it's committed.
func (r *Repository) ExecTxn(ctx context.Context, fn func(repo *Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new repository instance that is bound to the transaction.
		txnRepo := NewRepository(tx)
		// Execute the provided function with the transactional repository.
		return fn(txnRepo)
	})
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

// --- Create Methods ---

func (r *Repository) CreateContact(contact *model.Contact) error {
	return r.db.Create(contact).Error
}

// --- Get Methods ---

func (r *Repository) GetByID(id string) (*model.Contact, error) {
	var contact model.Contact
	if err := r.db.First(&contact, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *Repository) GetByQuery(query map[string]interface{}) (*model.Contact, error) {
	var contact model.Contact
	if err := r.db.Where(query).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// --- Update Methods ---

func (r *Repository) UpdateByID(contact *model.Contact) error {
	return r.db.Save(contact).Error
}

func (r *Repository) UpdateByQuery(query map[string]interface{}, values map[string]interface{}) error {
	return r.db.Model(&model.Contact{}).Where(query).Updates(values).Error
}

// --- Mappers ---

func ToContactsDbModel(req contracts.CreateContactRequest) (*model.Contact, error) {
	return &model.Contact{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		OwnerID:   req.OwnerID,
	}, nil
}

func (r *Repository) ToContactsApiResponse(contact *model.Contact) *contracts.CreateContactResponse {
	var owner User
	ownerName := "Unknown"

	if err := r.db.First(&owner, "id = ?", contact.OwnerID).Error; err == nil {
		ownerName = owner.FirstName + " " + owner.LastName
	} else {
		log.Printf("Could not find owner with ID %s: %v", contact.OwnerID, err)
	}

	return &contracts.CreateContactResponse{
		ID:        contact.ID,
		Name:      contact.FirstName + " " + contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		OwnerID:   contact.OwnerID,
		OwnerName: ownerName,
		CreatedAt: contact.CreatedAt,
	}
}
