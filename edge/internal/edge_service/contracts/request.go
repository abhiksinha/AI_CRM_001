package contracts

import "encoding/json"

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

// LoginRequest defines the JSON body for /v1/login.
type LoginRequest struct {
	Username string   `json:"username"`
	Password Password `json:"password"`
}

// LogoutRequest defines the JSON body for /v1/logout.
type LogoutRequest struct {
	APIToken string `json:"api_token"`
}

// SessionsRequest defines the JSON body for /v1/sessions.
type SessionsRequest struct {
	Username string   `json:"username"`
	Password Password `json:"password"`
}

// ExpireSessionRequest defines the JSON body for /v1/sessions/expire.
type ExpireSessionRequest struct {
	SessionID string `json:"session_id"`
}

// GetTokenRequest defines the JSON body for /v1/token.
type GetTokenRequest struct {
	APIKey string `json:"api_key"`
}
