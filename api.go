package telegrambot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type RemoteSource struct {
	FileIDOrURL string
}

func (RemoteSource) isFileSource() {}

func (s RemoteSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.FileIDOrURL)
}

func (s *RemoteSource) UnmarshalJSON(data []byte) error {
	var fileIDOrURL string
	if err := json.Unmarshal(data, &fileIDOrURL); err != nil {
		return err
	}
	s.FileIDOrURL = fileIDOrURL
	return nil
}

type sourcePrefix string

const (
	attachPrefix sourcePrefix = "attach://"
	filePrefix   sourcePrefix = "file://"
)

type LocalSource struct {
	prefix   sourcePrefix
	Filepath string
}

func (*LocalSource) isFileSource() {}

func filenameFromFilepath(path string) string {
	filename := filepath.Base(path)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func (s LocalSource) MarshalJSON() ([]byte, error) {
	filepath := s.Filepath
	if s.prefix == attachPrefix {
		filepath = filenameFromFilepath(s.Filepath)
	}
	return json.Marshal(string(s.prefix) + filepath)
}

func (s *LocalSource) UnmarshalJSON(data []byte) error {
	var rawFilepath string
	if err := json.Unmarshal(data, &rawFilepath); err != nil {
		return err
	}
	prefixInd := strings.Index(rawFilepath, "://")
	s.prefix = sourcePrefix(rawFilepath[:prefixInd])
	filepathInd := prefixInd + 3
	s.Filepath = rawFilepath[filepathInd:]
	return nil
}

type FileSource interface {
	isFileSource()
}

func (b *Bot) DownloadFile(
	ctx context.Context,
	fileID string,
	outDir string,
) (string, error) {
	method := "downloadFile"
	file, err := b.GetFile(ctx, GetFileParams{FileID: fileID})
	if err != nil {
		return "", err
	}
	if file.Path == "" {
		return "", errWithPrefix(method, ErrFileIsMissing)
	}
	filePath := file.Path
	if b.filesDir == "" {
		endpoint, err := buildEndpoint(b.baseURL+"/file", b.token, filePath)
		if err != nil {
			return "", errWithPrefix(method, err)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return "", errWithPrefix(method, err)
		}
		respBody, statusCode, err := do(req, b.httpClient)
		if err != nil {
			return "", errWithPrefix(method, err)
		}
		if statusCode > 199 && statusCode < 300 {
			filename := filepath.Base(filePath)
			filePath = filepath.Join(outDir, filename)
			if err := os.WriteFile(filePath, respBody, 0644); err != nil {
				return "", errWithPrefix(method, err)
			}
		} else {
			if err := parseResponse(respBody, statusCode, nil); err != nil {
				return "", errWithPrefix(method, err)
			} else {
				return "", errWithPrefix(method, fmt.Errorf("unexpected error (status %d)", statusCode))
			}
		}
	}
	return filePath, nil
}
