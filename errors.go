package telegrambot

import "fmt"

const errPrefix string = "telegrambot"

func errWithPrefix(method string, err error) error {
    prefix := fmt.Sprintf("%s: %s", errPrefix, method)
    return fmt.Errorf("%s: %w",prefix,err)
}
