package contracts

import (
	"encoding/json"
)

// Password is a custom type to prevent accidental logging of plain text passwords.
type Password string

func (p Password) String() string               { return "[REDACTED]" }
func (p Password) MarshalJSON() ([]byte, error) { return json.Marshal(string(p)) }
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
	Password  Password `json:"password"`
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
	Password Password `json:"password"`
}

// CreateApiKeyRequest defines the JSON body for creating an API key.
type CreateApiKeyRequest struct {
	UserID string `json:"user_id"`
}

// MatchApiKeyRequest defines the JSON body for matching an API key.
type MatchApiKeyRequest struct {
	ApiKey string `json:"api_key"`
}

// ExpireApiKeyRequest defines the JSON body for expiring an API key.
type ExpireApiKeyRequest struct {
	ApiKey string `json:"api_key"`
}
