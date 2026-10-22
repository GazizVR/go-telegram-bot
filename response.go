package telegrambot

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

func parseResponse(
	respBody []byte,
	statusCode int,
	result any,
) error {
	isSuccessful := statusCode >= 200 && statusCode < 300
	var resp response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		if !isSuccessful {
			errText := "unexpected error"
			httpStatusText := http.StatusText(statusCode)
			if httpStatusText != "" {
				errText = httpStatusText
			}
			return fmt.Errorf("%s (status %d)", errText, statusCode)
		}
		return fmt.Errorf("unmarshal respBody: %w", err)
	}
	if !resp.OK {
		errText := "unexpected error"
		if resp.Description != "" {
			errText = resp.Description
		} else {
			if !isSuccessful {
				httpStatusText := http.StatusText(statusCode)
				if httpStatusText != "" {
					errText = httpStatusText
				}
			}
		}
		if resp.ErrorCode != 0 {
			errText = errText + fmt.Sprintf(" (status %d)", resp.ErrorCode)
		} else {
			if !isSuccessful {
				errText = errText + fmt.Sprintf(" (status %d)", statusCode)
			}
		}
		return fmt.Errorf("%s", errText)
	}
	if !isSuccessful {
		errText := "unexpected error"
		httpStatusText := http.StatusText(statusCode)
		if httpStatusText != "" {
			errText = httpStatusText
		}
		return fmt.Errorf("%s (status %d)", errText, statusCode)
	}
	if result != nil && len(resp.Result) > 0 {
		if err := json.Unmarshal(resp.Result, result); err != nil {
			return fmt.Errorf("unmarshal result: %w", err)
		}
	}
	return nil
}
