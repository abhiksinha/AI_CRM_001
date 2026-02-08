package service

import (
	"context"
	"edge/packages/httpRequest"
	"edge/packages/public_response"
	"net/url"
)

type listUsersResponse struct {
	Data       []userResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type verifyPasswordResponse struct {
	IsValid bool `json:"is_valid"`
}

type matchApiKeyResponse struct {
	IsValid bool   `json:"is_valid"`
	UserID  string `json:"user_id"`
}

func (s *AuthService) getUserIDByEmail(ctx context.Context, email string) (string, error) {
	q := url.Values{}
	q.Set("page", "1")
	q.Set("page_size", "1")
	q.Set("email", email)

	var payload listUsersResponse
	endpoint := "/api/v1/users?" + q.Encode()
	if err := httpRequest.MakeRequest(ctx, s.userServiceBase, endpoint, nil, &payload); err != nil {
		return "", err
	}
	if len(payload.Data) == 0 {
		return "", public_response.ErrNotFound
	}
	return payload.Data[0].ID, nil
}

func (s *AuthService) verifyPassword(ctx context.Context, userID, password string) (bool, error) {
	var payload verifyPasswordResponse
	endpoint := "/api/v1/users/verify-password"
	req := struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}{
		ID:       userID,
		Password: password,
	}
	if err := httpRequest.MakeRequest(ctx, s.userServiceBase, endpoint, req, &payload); err != nil {
		return false, err
	}
	return payload.IsValid, nil
}

func (s *AuthService) matchApiKey(ctx context.Context, apiKey string) (string, error) {
	endpoint := "/api/v1/users/api-keys/match"
	req := struct {
		ApiKey string `json:"api_key"`
	}{
		ApiKey: apiKey,
	}
	var payload matchApiKeyResponse
	if err := httpRequest.MakeRequest(ctx, s.userServiceBase, endpoint, req, &payload); err != nil {
		return "", err
	}
	if !payload.IsValid || payload.UserID == "" {
		return "", public_response.ErrUnauthorized
	}
	return payload.UserID, nil
}
