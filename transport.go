package telegrambot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func do(
	ctx context.Context,
	baseURL string,
	token string,
	method string,
	body any,
	httpClient http.Client,
) (respBody []byte, statusCode int, err error) {
	endpoint, err := url.JoinPath(baseURL, "bot"+token, method)
	if err != nil {
		return nil, 0, fmt.Errorf("create endpoint: %w", err)
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, 0, fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read response: %w", err)
	}
	return respBody, resp.StatusCode, nil
}
