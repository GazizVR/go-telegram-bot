package telegrambot

import "fmt"

const errPrefix string = "telegrambot"

func errorPrefix(method string) string {
	return fmt.Sprintf("%s: %s", errPrefix, method)
}
