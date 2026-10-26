package telegrambot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func WithFilesDir(filesDir string) Option {
	return func(b *Bot) error {
		dir, err := os.Open(filesDir)
		if err != nil {
			return fmt.Errorf("%s: open files dir by provided path: %s, error: %w", errPrefix, filesDir, err)
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

func fileSourcesToPaths(
	filesDir string,
	fileSources ...FileSource,
) (map[string]string, error) {
	var filePaths map[string]string
	if filesDir != "" {
		for _, s := range fileSources {
			local, ok := s.(*LocalSource)
			if ok {
				filename := filepath.Base(local.Filepath)
				filePath := filepath.Join(filesDir, filename)
				srcFile, err := os.Open(local.Filepath)
				if err != nil {
					return nil, fmt.Errorf("open %s file: %w", filename, err)
				}
				defer srcFile.Close()
				dstFile, err := os.Create(filePath)
				if err != nil {
					return nil, fmt.Errorf("create %s file copy: %w", filename, err)
				}
				if _, err := io.Copy(dstFile, srcFile); err != nil {
					return nil, fmt.Errorf("write data into %s file: %w", filename, err)
				}
				if err := dstFile.Close(); err != nil {
					return nil, fmt.Errorf("close %s file: %w", filename, err)
				}
				local.prefix = filePrefix
				local.Filepath = filePath
			}
		}
	} else {
		filePaths = make(map[string]string)
		for _, s := range fileSources {
			local, ok := s.(*LocalSource)
			if ok {
				local.prefix = attachPrefix
				filePaths[filenameFromFilepath(local.Filepath)] = local.Filepath
			}
		}
	}
	return filePaths, nil
}

func (b *Bot) rawCall(
	ctx context.Context,
	method string,
	body any,
	result any,
	fileSources ...FileSource,
) error {
	var (
		respBody   []byte
		statusCode int
		err        error
	)
	filePaths, err := fileSourcesToPaths(b.filesDir, fileSources...)
	if err != nil {
		return err
	}
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

func (b *Bot) call(
	ctx context.Context,
	method string,
	body any,
	result any,
	fileSources ...FileSource,
) error {
	if err := b.rawCall(ctx, method, body, result, fileSources...); err != nil {
		prefix := fmt.Sprintf("%s: %s", errPrefix, method)
		return fmt.Errorf("%s: %w", prefix, err)
	}
	return nil
}
