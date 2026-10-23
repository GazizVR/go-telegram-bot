package telegrambot

type Message struct {
    ID int64 `json:"message_id"`
}

type Update struct {
	ID int64 `json:"update_id"`
    Message *Message `json:"message"`
}
