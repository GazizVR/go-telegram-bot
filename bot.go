package telegrambot

import (
	"net/http"
)

type Bot struct {
	baseURL    string
	token      string
	httpClient http.Client
}

const defaultBaseURL = "https://api.telegram.org"

func NewBot(token string) *Bot {
	return &Bot{
		baseURL:    defaultBaseURL,
		token:      token,
		httpClient: *http.DefaultClient,
	}
}
