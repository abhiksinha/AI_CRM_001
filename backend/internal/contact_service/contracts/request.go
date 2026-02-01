package contracts

// ListContactsRequest defines the query parameters for listing contacts.
// These will be read from the URL query string, not a JSON body.
type ListContactsRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	SortBy   string `form:"sort_by"` // e.g., "created_at_desc"
	Query    string `form:"query"`   // For full-text search
}

// CreateContactRequest defines the expected JSON body for a create contact request.
type CreateContactRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	OwnerID   string `json:"owner_id"`
}

// UpdateContactRequest defines the expected JSON body for an update contact request.
// Pointers are used to distinguish between a field that is intentionally set to an empty value
// and a field that is not being updated at all.
type UpdateContactRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	OwnerID   *string `json:"owner_id"`
}

// AddNoteRequest defines the expected JSON body for adding a note.
type AddNoteRequest struct {
	Content string `json:"content"`
}
