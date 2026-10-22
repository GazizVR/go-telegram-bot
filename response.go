package telegrambot

import (
	"encoding/json"
	"fmt"
)

type response struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

func parseResponse(
	respBody []byte,
	result any,
) error {
	var resp response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("unmarshal respBody: %w", err)
	}
	return nil
}
