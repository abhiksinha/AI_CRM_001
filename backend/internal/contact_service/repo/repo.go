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
func (r *Repository) ExecTxn(ctx context.Context, fn func(repo *Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txnRepo := NewRepository(tx)
		return fn(txnRepo)
	})
}

// User is a minimal struct to fetch user data for the response.
type User struct {
	ID        string `gorm:"type:varchar(36);primary_key"`
	FirstName string
	LastName  string
}

func (u *User) TableName() string { return "users" }

// --- Create Methods ---

func (r *Repository) CreateContact(contact *model.Contact) error {
	return r.db.Create(contact).Error
}

// --- Get Methods ---

func (r *Repository) GetByIDAndOwner(id, ownerID string) (*model.Contact, error) {
	var contact model.Contact
	err := r.db.First(&contact, "id = ? AND owner_id = ?", id, ownerID).Error
	return &contact, err
}

func (r *Repository) ListByOwner(ownerID string, req contracts.ListContactsRequest) ([]model.Contact, int64, error) {
	var contacts []model.Contact
	var totalCount int64

	query := r.db.Model(&model.Contact{}).Where("owner_id = ?", ownerID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	err := query.Find(&contacts).Error
	return contacts, totalCount, err
}

func (r *Repository) GetByQuery(query map[string]interface{}) (*model.Contact, error) {
	var contact model.Contact
	err := r.db.Where(query).First(&contact).Error
	if contact.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &contact, err
}

// --- Update Methods ---

func (r *Repository) UpdateContact(contact *model.Contact) error {
	return r.db.Save(contact).Error
}

// --- Delete Methods ---

func (r *Repository) DeleteByIDAndOwner(id, ownerID string) (int64, error) {
	result := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Contact{})
	return result.RowsAffected, result.Error
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

func (r *Repository) ToContactsApiResponse(contact *model.Contact) *contracts.ContactResponse {
	var owner User
	ownerName := "Unknown"

	if err := r.db.First(&owner, "id = ?", contact.OwnerID).Error; err == nil {
		ownerName = owner.FirstName + " " + owner.LastName
	} else {
		log.Printf("Could not find owner with ID %s: %v", contact.OwnerID, err)
	}

	return &contracts.ContactResponse{
		ID:        contact.ID,
		Name:      contact.FirstName + " " + contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		OwnerID:   contact.OwnerID,
		OwnerName: ownerName,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}
