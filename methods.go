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
		return nil, err
	}
	return updates, nil
}

type SendMessageParams struct {
	ChatID      int64            `json:"chat_id"`
	Text        string           `json:"text"`
	ReplyParams *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup ReplyMarkup      `json:"reply_markup,omitempty"`
}

const sendMessageMethod string = "sendMessage"

func (b *Bot) SendMessage(
	ctx context.Context,
	params SendMessageParams,
) (*Message, error) {
	var msg Message
	if err := b.call(ctx, sendMessageMethod, params, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

type AnswerCallbackQueryParams struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

const answerCallbackQueryMethod string = "answerCallbackQuery"

func (b *Bot) AnswerCallbackQuery(
	ctx context.Context,
	params AnswerCallbackQueryParams,
) error {
	return b.call(ctx, answerCallbackQueryMethod, params, nil)
}

type EditMessageTextParams struct {
	ChatID      int64                 `json:"chat_id"`
	MessageID   int64                 `json:"message_id"`
	Text        string                `json:"text"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

const editMessageTextMethod string = "editMessageText"

func (b *Bot) EditMessageText(
	ctx context.Context,
	params EditMessageTextParams,
) (*Message, error) {
	var msg Message
	if err := b.call(ctx, editMessageTextMethod, params, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

type EditMessageReplyMarkupParams struct {
	ChatID      int64                 `json:"chat_id"`
	MessageID   int64                 `json:"message_id"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

const editMessageReplyMarkupMethod string = "editMessageReplyMarkup"

func (b *Bot) EditMessageReplyMarkup(
	ctx context.Context,
	params EditMessageReplyMarkupParams,
) (*Message, error) {
	var msg Message
	if err := b.call(ctx, editMessageReplyMarkupMethod, params, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
