package contracts

// LoginResponse is the response for /v1/login.
type LoginResponse struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	APIToken  string `json:"api_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// SessionsResponse is the response for /v1/sessions.
type SessionsResponse struct {
	SessionIDs []string `json:"session_ids"`
}

// GetTokenResponse is the response for /v1/token.
type GetTokenResponse struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	ExpiresIn int64  `json:"expires_in"`
}
