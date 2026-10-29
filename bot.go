package telegrambot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Bot struct {
	baseURL    string
	token      string
	httpClient http.Client
	filesDir   string
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

const errPrefix string = "telegrambot"

func errWithPrefix(method string, err error) error {
	return fmt.Errorf("%s: %s: %w", errPrefix, method, err)
}

func WithFilesDir(filesDir string) Option {
	return func(b *Bot) error {
		if strings.TrimSpace(filesDir) == "" {
			errText := fmt.Sprintf("%s: specified filesDir must not by empty", errPrefix)
			return errors.New(errText)
		}
		dir, err := os.Open(filesDir)
		if err != nil {
			return fmt.Errorf("%s: open filesDir: %s, error: %w", errPrefix, filesDir, err)
		}
		dir.Close()
		b.filesDir = filesDir
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
	fileSources ...FileSource,
) error {
	var (
		req *http.Request
		err error
	)
	endpoint, err := buildEndpoint(b.baseURL, b.token, method)
	if err != nil {
		return err
	}
	filePaths, err := fileSourcesToPaths(b.filesDir, fileSources...)
	if err != nil {
		return err
	}
	if len(filePaths) > 0 {
		req, err = buildMultipartRequest(ctx, body, filePaths, endpoint)
		if err != nil {
			return err
		}
	} else {
		req, err = buildJSONRequest(ctx, endpoint, body)
		if err != nil {
			return err
		}
	}
	respBody, statusCode, err := do(req, b.httpClient)
	if err != nil {
		return err
	}
	return parseResponse(respBody, statusCode, result)
}

func (b *Bot) call(
	ctx context.Context,
	method string,
	body any,
	result any,
	fileSources ...FileSource,
) error {
	if err := b.rawCall(ctx, method, body, result, fileSources...); err != nil {
		return errWithPrefix(method, err)
	}
	return nil
}
