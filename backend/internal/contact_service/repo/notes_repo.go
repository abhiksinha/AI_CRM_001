package repo

import (
	"CRM/internal/contact_service/model"
	"gorm.io/gorm"
)

// NotesRepository provides access to the database for the notes service.
type NotesRepository struct {
	db *gorm.DB
}

// NewNotesRepository creates a new notes repository.
func NewNotesRepository(db *gorm.DB) *NotesRepository {
	return &NotesRepository{db: db}
}

// CreateNote inserts a new note into the database.
func (r *NotesRepository) CreateNote(note *model.Note) error {
	return r.db.Create(note).Error
}

// ListByContactID retrieves all notes for a given contact.
func (r *NotesRepository) ListByContactID(contactID string) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.Where("contact_id = ?", contactID).Order("created_at desc").Find(&notes).Error
	return notes, err
}
