package contracts

// CreateContactRequest defines the expected JSON body for a create contact request.
type CreateContactRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	OwnerID   string `json:"owner_id"`
}
