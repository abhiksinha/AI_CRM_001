package service

import (
	"context"
	"crypto/rand"
	"edge/internal/edge_service/contracts"
	"edge/packages/public_response"
	edgeredis "edge/packages/redis"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	maxSessions = 3
	sessionTTL  = 30 * time.Minute
)

type AuthService struct {
	redis           *goredis.Client
	userServiceBase string
	httpClient      *http.Client
	logger          *zap.Logger
}

func NewAuthService(redisClient *goredis.Client, userServiceBase string, logger *zap.Logger) *AuthService {
	return &AuthService{
		redis:           redisClient,
		userServiceBase: strings.TrimRight(userServiceBase, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}

func (s *AuthService) AllowRate(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	return edgeredis.AllowRate(ctx, s.redis, key, limit, window)
}

func (s *AuthService) Login(ctx context.Context, req contracts.LoginRequest) (*contracts.LoginResponse, error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(string(req.Password)) == "" {
		return nil, public_response.ErrValidation
	}

	userID, err := s.resolveUserID(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	isValid, err := s.verifyPassword(ctx, userID, string(req.Password))
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, public_response.ErrUnauthorized
	}

	if err := s.enforceSessionLimit(ctx, userID); err != nil {
		return nil, err
	}

	sessionID, err := randomToken(32)
	if err != nil {
		return nil, err
	}
	apiToken, err := randomToken(32)
	if err != nil {
		return nil, err
	}

	if err := edgeredis.PersistSession(ctx, s.redis, userID, sessionID, apiToken, sessionTTL); err != nil {
		s.logger.Error("failed to persist session", zap.Error(err))
		return nil, err
	}

	return &contracts.LoginResponse{
		UserID:    userID,
		SessionID: sessionID,
		APIToken:  apiToken,
		ExpiresIn: int64(sessionTTL.Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req contracts.LogoutRequest) error {
	if strings.TrimSpace(req.APIToken) == "" {
		return public_response.ErrValidation
	}

	tokenKey := fmt.Sprintf("auth:token:%s", req.APIToken)
	data, err := s.redis.HGetAll(ctx, tokenKey).Result()
	if err != nil {
		s.logger.Error("failed to fetch token data", zap.Error(err))
		return err
	}
	if len(data) == 0 {
		return public_response.ErrUnauthorized
	}

	userID := data["user_id"]
	sessionID := data["session_id"]
	if userID == "" || sessionID == "" {
		return public_response.ErrUnauthorized
	}

	if err := edgeredis.ClearSession(ctx, s.redis, userID, sessionID, req.APIToken); err != nil {
		s.logger.Error("failed to clear session", zap.Error(err))
		return err
	}

	return nil
}

func (s *AuthService) ListSessions(ctx context.Context, req contracts.SessionsRequest) (*contracts.SessionsResponse, error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(string(req.Password)) == "" {
		return nil, public_response.ErrValidation
	}

	userID, err := s.resolveUserID(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	isValid, err := s.verifyPassword(ctx, userID, string(req.Password))
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, public_response.ErrUnauthorized
	}

	sessions, err := edgeredis.ActiveSessions(ctx, s.redis, userID)
	if err != nil {
		s.logger.Error("failed to list sessions", zap.Error(err))
		return nil, err
	}

	return &contracts.SessionsResponse{SessionIDs: sessions}, nil
}

func (s *AuthService) ExpireSession(ctx context.Context, req contracts.ExpireSessionRequest) error {
	if strings.TrimSpace(req.SessionID) == "" {
		return public_response.ErrValidation
	}

	if err := edgeredis.ExpireSessionByID(ctx, s.redis, req.SessionID); err != nil {
		s.logger.Error("failed to expire session", zap.Error(err))
		return err
	}
	return nil
}

func (s *AuthService) GetToken(ctx context.Context, req contracts.GetTokenRequest) (*contracts.GetTokenResponse, error) {
	if strings.TrimSpace(req.APIKey) == "" {
		return nil, public_response.ErrValidation
	}

	userID, err := s.matchApiKey(ctx, req.APIKey)
	if err != nil {
		return nil, err
	}

	if err := s.enforceSessionLimit(ctx, userID); err != nil {
		return nil, err
	}

	sessionID, err := randomToken(32)
	if err != nil {
		return nil, err
	}
	apiToken, err := randomToken(32)
	if err != nil {
		return nil, err
	}

	if err := edgeredis.PersistSession(ctx, s.redis, userID, sessionID, apiToken, sessionTTL); err != nil {
		s.logger.Error("failed to persist session", zap.Error(err))
		return nil, err
	}

	return &contracts.GetTokenResponse{
		Token:     apiToken,
		SessionID: sessionID,
		ExpiresIn: int64(sessionTTL.Seconds()),
	}, nil
}

func (s *AuthService) resolveUserID(ctx context.Context, username string) (string, error) {
	username = strings.TrimSpace(username)
	if strings.Contains(username, "@") {
		return s.getUserIDByEmail(ctx, username)
	}
	return username, nil
}

func (s *AuthService) enforceSessionLimit(ctx context.Context, userID string) error {
	sessionSetKey := fmt.Sprintf("auth:sessions:%s", userID)
	members, err := s.redis.SMembers(ctx, sessionSetKey).Result()
	if err != nil {
		s.logger.Error("failed to read session set", zap.Error(err))
		return err
	}

	for _, sessionID := range members {
		sessionKey := fmt.Sprintf("auth:session:%s", sessionID)
		exists, err := s.redis.Exists(ctx, sessionKey).Result()
		if err != nil {
			s.logger.Error("failed to check session key", zap.Error(err))
			return err
		}
		if exists == 0 {
			_ = s.redis.SRem(ctx, sessionSetKey, sessionID).Err()
		}
	}

	count, err := s.redis.SCard(ctx, sessionSetKey).Result()
	if err != nil {
		s.logger.Error("failed to count sessions", zap.Error(err))
		return err
	}
	if count >= maxSessions {
		return public_response.ErrSessionLimit
	}
	return nil
}

func randomToken(bytesLen int) (string, error) {
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
