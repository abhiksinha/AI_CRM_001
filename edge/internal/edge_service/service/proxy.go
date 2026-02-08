package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"edge/packages/public_response"

	"github.com/go-chi/chi/v5"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ProxyService struct {
	redis  *goredis.Client
	client *http.Client
	logger *zap.Logger
}

func NewProxyService(redisClient *goredis.Client, logger *zap.Logger) *ProxyService {
	return &ProxyService{
		redis: redisClient,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (s *ProxyService) VerifyToken(ctx context.Context, token string) (string, error) {
	tokenKey := fmt.Sprintf("auth:token:%s", token)
	data, err := s.redis.HGetAll(ctx, tokenKey).Result()
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", public_response.ErrUnauthorized
	}
	userID := data["user_id"]
	if userID == "" {
		return "", public_response.ErrUnauthorized
	}
	return userID, nil
}

func (s *ProxyService) Forward(w http.ResponseWriter, r *http.Request, baseURL, backendPath, userID string) error {
	resolvedPath, err := fillPathParams(backendPath, r)
	if err != nil {
		return err
	}

	target := strings.TrimRight(baseURL, "/") + resolvedPath
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequestWithContext(r.Context(), r.Method, target, r.Body)
	if err != nil {
		return err
	}

	copyHeaders(req.Header, r.Header)
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	return err
}

func fillPathParams(path string, r *http.Request) (string, error) {
	resolved := path
	start := strings.Index(resolved, "{")
	for start != -1 {
		end := strings.Index(resolved[start:], "}")
		if end == -1 {
			break
		}
		end += start
		param := resolved[start+1 : end]
		val := chi.URLParam(r, param)
		if val == "" {
			return "", fmt.Errorf("missing url param: %s", param)
		}
		resolved = resolved[:start] + val + resolved[end+1:]
		start = strings.Index(resolved, "{")
	}
	return resolved, nil
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		if shouldSkipHeader(k) {
			continue
		}
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func shouldSkipHeader(key string) bool {
	switch http.CanonicalHeaderKey(key) {
	case "Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization", "Te", "Trailer", "Transfer-Encoding", "Upgrade", "Content-Length":
		return true
	default:
		return false
	}
}
