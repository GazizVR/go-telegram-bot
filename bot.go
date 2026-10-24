package telegrambot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Bot struct {
	baseURL    string
	token      string
	httpClient http.Client
}

type Option func(*Bot) error

func WithBaseURL(u string) Option {
	return func(b *Bot) error {
		url, err := url.Parse(u)
		if err != nil {
			return err
		}
		if url.Host == "" {
			return errors.New("url host must not be empty")
		}
		if url.Scheme != "https" && url.Scheme != "http" {
			return errors.New("unsupported url scheme, use http or https")
		}
		b.baseURL = url.String()
		return nil
	}
}

func WithHTTPCLient(client http.Client) Option {
	return func(b *Bot) error {
		b.httpClient = client
		return nil
	}
}

const defaultBaseURL = "https://api.telegram.org"

func NewBot(token string, opts ...Option) (*Bot, error) {
	bot := &Bot{
		baseURL:    defaultBaseURL,
		token:      token,
		httpClient: *http.DefaultClient,
	}
	for _, op := range opts {
		if err := op(bot); err != nil {
			return nil, err
		}
	}
	return bot, nil
}

func (b *Bot) rawCall(
	ctx context.Context,
	method string,
	body any,
	result any,
	filePaths map[string]string,
) error {
	var (
		respBody   []byte
		statusCode int
		err        error
	)
	if len(filePaths) > 0 {
		respBody, statusCode, err = doMultipart(ctx, body, filePaths, b.baseURL, b.token, method, b.httpClient)
		if err != nil {
			return err
		}
	} else {
		respBody, statusCode, err = do(ctx, b.baseURL, b.token, method, body, b.httpClient)
		if err != nil {
			return err
		}
	}
	return parseResponse(respBody, statusCode, result)
}

const errPrefix string = "telegrambot"

func (b *Bot) call(
	ctx context.Context,
	method string,
	body any,
	result any,
	filePaths map[string]string,
) error {
	if err := b.rawCall(ctx, method, body, result, filePaths); err != nil {
		prefix := fmt.Sprintf("%s: %s", errPrefix, method)
		return fmt.Errorf("%s: %w", prefix, err)
	}
	return nil
}
