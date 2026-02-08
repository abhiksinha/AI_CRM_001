package contracts

// UserResponse is the standard response for a single user.
type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// ListUsersResponse is the response for listing multiple users.
type ListUsersResponse struct {
	Data       []UserResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
}

// ApiKeyResponse is the standard response for API key operations.
type ApiKeyResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
}

// ApiKeyMatchResponse is the response for matching an API key.
type ApiKeyMatchResponse struct {
	IsValid bool   `json:"is_valid"`
	UserID  string `json:"user_id"`
}
