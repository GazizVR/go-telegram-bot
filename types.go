package telegrambot

import "encoding/json"

type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	URL              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

type Chat struct {
	ID               int64  `json:"id"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	Username         string `json:"username"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	IsForum          bool   `json:"is_forum"`
	IsDirectMessages bool   `json:"is_direct_messages"`
}

type Message struct {
	ID                 int64               `json:"message_id"`
	Chat               Chat                `json:"chat"`
	Text               string              `json:"text,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

func (Message) maybeInaccessibleMessage() {}

type InaccessibleMessage struct {
	Chat Chat  `json:"chat"`
	ID   int64 `json:"message_id"`
}

func (InaccessibleMessage) maybeInaccessibleMessage() {}

type MaybeInaccessibleMessage interface {
	maybeInaccessibleMessage()
}

type CallbackQuery struct {
	ID              string                   `json:"id"`
	Message         MaybeInaccessibleMessage `json:"message,omitempty"`
	Data            string                   `json:"data,omitempty"`
	InlineMessageID string                   `json:"inline_message_id,omitempty"`
	ChatInstance    string                   `json:"chat_instance,omitempty"`
	GameShortName   string                   `json:"game_short_name,omitempty"`
}

func (q *CallbackQuery) UnmarshalJSON(data []byte) error {
	type Alias CallbackQuery
	aux := &struct {
		*Alias
		RawMessage json.RawMessage `json:"message,omitempty"`
	}{
		Alias: (*Alias)(q),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	type maybeInaccessibleMessage struct {
		Date int64 `json:"date"`
	}
	var msg maybeInaccessibleMessage
	if err := json.Unmarshal(aux.RawMessage, &msg); err != nil {
		return err
	}
	if msg.Date == 0 {
		var newMsg InaccessibleMessage
		if err := json.Unmarshal(aux.RawMessage, &newMsg); err != nil {
			return err
		}
		q.Message = newMsg
	} else {
		var newMsg Message
		if err := json.Unmarshal(aux.RawMessage, &newMsg); err != nil {
			return err
		}
		q.Message = newMsg
	}
	return nil
}

type Update struct {
	ID            int64          `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	Style        string `json:"style,omitempty"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

type InlineKeyboardMarkup struct {
	Keyboard   [][]InlineKeyboardButton `json:"inline_keyboard"`
	ForceReply bool                     `json:"force_reply,omitempty"`
}

func (InlineKeyboardMarkup) isReplyMarkup() {}

type ReplyMarkup interface {
	isReplyMarkup()
}

type ReplyParameters struct {
	MessageID *int64 `json:"message_id,omitempty"`
}

type MediaType string

const (
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
)

type InputMedia struct {
	Type  MediaType  `json:"type"`
	Media FileSource `json:"media"`
}
