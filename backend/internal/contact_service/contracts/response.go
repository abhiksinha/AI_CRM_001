package contracts

// ContactResponse is the standard response for a single contact.
// It is used for Get, Create, and Update responses.
type ContactResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	OwnerID   string `json:"owner_id"`
	OwnerName string `json:"owner_name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// ListContactsResponse is the response for listing multiple contacts.
type ListContactsResponse struct {
	Data       []ContactResponse `json:"data"`
	TotalCount int64             `json:"total_count"`
}

// NoteResponse is the standard response for a single note.
type NoteResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	AuthorID  string `json:"author_id"`
	CreatedAt int64  `json:"created_at"`
}

// ListNotesResponse is the response for listing multiple notes.
type ListNotesResponse struct {
	Data       []NoteResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
}
