package contracts

import (
	"encoding/json"
	"fmt"
)

// Password is a custom type to prevent accidental logging of plain text passwords.
type Password string

// String implements the fmt.Stringer interface.
// This is automatically called by most logging libraries (like zap) and fmt.Println.
func (p Password) String() string {
	return "[REDACTED]"
}

// MarshalJSON implements the json.Marshaler interface.
// This ensures the password is sent as a plain string in JSON requests.
func (p Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This allows the password to be read as a plain string from JSON requests.
func (p *Password) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*p = Password(s)
	return nil
}

// ListUsersRequest defines the query parameters for listing users.
type ListUsersRequest struct {
	Page     int    `schema:"page"`
	PageSize int    `schema:"page_size"`
	Role     string `schema:"role"`
}

// CreateUserRequest defines the JSON body for creating a new user.
type CreateUserRequest struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Password  Password `json:"password"` // Use the new Password type
	Role      string   `json:"role"`
}

// UpdateUserRequest defines the JSON body for updating a user.
type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      *string `json:"role"`
	IsActive  *bool   `json:"is_active"`
}

// VerifyPasswordRequest defines the JSON body for the password verification endpoint.
type VerifyPasswordRequest struct {
	ID       string   `json:"id"`
	Password Password `json:"password"` // Use the new Password type
}
