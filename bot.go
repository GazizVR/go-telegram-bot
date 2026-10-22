package telegrambot

import (
	"context"
	"net/http"
)

type Bot struct {
	baseURL    string
	token      string
	httpClient http.Client
}

const defaultBaseURL = "https://api.telegram.org"

func NewBot(token string) *Bot {
	return &Bot{
		baseURL:    defaultBaseURL,
		token:      token,
		httpClient: *http.DefaultClient,
	}
}

func (b *Bot) call(
	ctx context.Context,
	method string,
	body any,
	result any,
) error {
	respBody, statusCode, err := do(ctx, b.baseURL, b.token, method, body, b.httpClient)
	if err != nil {
		return err
	}
	return parseResponse(respBody, statusCode, result)
}
