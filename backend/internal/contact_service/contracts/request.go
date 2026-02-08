package contracts

// ListContactsRequest defines the query parameters for listing contacts.
type ListContactsRequest struct {
	Page     int    `schema:"page"`
	PageSize int    `schema:"page_size"`
	SortBy   string `schema:"sort_by"`
	// A generic 'query' parameter has been removed in favor of specific filters if needed later.
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
