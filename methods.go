package telegrambot

import "context"

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
		return nil, errWithPrefix(getUpdatesMethod, err)
	}
	return updates, nil
}

type SendMessageParams struct {
	ChatID      int64            `json:"chat_id"`
	Text        string           `json:"text"`
	ReplyParams *ReplyParameters `json:"reply_parameters,omitempty"`
}

const sendMessageMethod string = "sendMessage"

func (b *Bot) SendMessage(
	ctx context.Context,
	params SendMessageParams,
) (*Message, error) {
	var msg Message
	if err := b.call(ctx, sendMessageMethod, params, &msg); err != nil {
		return nil, errWithPrefix(sendMessageMethod, err)
	}
	return &msg, nil
}
