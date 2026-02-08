package httpRequest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultTimeout = 5 * time.Second

// HTTPError preserves upstream status and body for pass-through.
type HTTPError struct {
	StatusCode int
	Body       []byte
}

func (e *HTTPError) Error() string {
	return http.StatusText(e.StatusCode)
}

// MakeRequest sends a JSON request to baseURL+endpoint and decodes the JSON response into responsePayload.
// If requestPayload is nil, it sends a GET; otherwise it sends a POST.
func MakeRequest(ctx context.Context, baseURL, endpoint string, requestPayload interface{}, responsePayload interface{}) error {
	client := &http.Client{Timeout: defaultTimeout}

	url := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(endpoint, "/")
	method := http.MethodGet
	if requestPayload != nil {
		method = http.MethodPost
	}

	var body io.Reader
	if requestPayload != nil {
		b, err := json.Marshal(requestPayload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	if requestPayload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		b, _ := io.ReadAll(resp.Body)
		return &HTTPError{StatusCode: resp.StatusCode, Body: b}
	}

	if responsePayload == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(responsePayload); err != nil {
		return err
	}
	return nil
}
