package telegrambot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func buildEndpoint(
	baseURL string,
	token string,
	method string,
) (string, error) {
	endpoint, err := url.JoinPath(baseURL, "bot"+token, method)
	if err != nil {
		return "", fmt.Errorf("create endpoint: %w", err)
	}
	return endpoint, nil
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
				if filePath != local.Filepath {
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

func do(
	req *http.Request,
	httpClient http.Client,
) (respBody []byte, statusCode int, err error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read response: %w", err)
	}
	return respBody, resp.StatusCode, nil
}

func buildJSONRequest(
	ctx context.Context,
	endpoint string,
	body any,
) (*http.Request, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func buildMultipartRequest(
	ctx context.Context,
	reqBody any,
	filePaths map[string]string,
	endpoint string,
) (*http.Request, error) {
	rawBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("encode reqBody: %w", err)
	}
	var rawFields map[string]json.RawMessage
	if err := json.Unmarshal(rawBody, &rawFields); err != nil {
		return nil, fmt.Errorf("decode rawBody: %w", err)
	}
	fields := make(map[string]string)
	for k, v := range rawFields {
		var str string
		if err := json.Unmarshal(v, &str); err != nil {
			fields[k] = fmt.Sprintf("%v", v)
		} else {
			fields[k] = str
		}
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("create %s field: %w", k, err)
		}
	}
	for k, v := range filePaths {
		file, err := os.Open(v)
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}
		filename := filepath.Base(file.Name())
		formFile, err := writer.CreateFormFile(k, filename)
		if err != nil {
			return nil, fmt.Errorf("create %s form file: %w", k, err)
		}
		if _, err := io.Copy(formFile, file); err != nil {
			return nil, fmt.Errorf("write data into %s file: %w", k, err)
		}
		_ = file.Close()
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}
