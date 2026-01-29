package contracts

import "time"

// CreateContactResponse defines the JSON body for a successful contact creation response.
type CreateContactResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	OwnerID   string    `json:"owner_id"`
	OwnerName string    `json:"owner_name"`
	CreatedAt time.Time `json:"created_at"`
}
