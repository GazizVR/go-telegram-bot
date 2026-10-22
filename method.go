package telegrambot

import (
	"context"
	"fmt"
)

type GetUpdatesParams struct {
	Offset         int      `json:"offset,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	Timeout        int      `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

const getUpdatesMethod string = "getUpdates"

func (b *Bot) GetUpdates(
	ctx context.Context,
	params GetUpdatesParams,
) ([]Update, error) {
	var updates []Update
	if err := b.call(ctx, getUpdatesMethod, params, &updates); err != nil {
		return nil, fmt.Errorf("%s: %w", errorPrefix(getUpdatesMethod), err)
	}
	return updates, nil
}
