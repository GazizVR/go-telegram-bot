package telegrambot

import "strings"

type InternalSource struct {
	FileID string
}

func (InternalSource) isFileSource() {}

func (s InternalSource) MarsahlJSON() ([]byte, error) {
	return []byte(s.FileID), nil
}

func (s *InternalSource) UnmarshalJSON(data []byte) error {
	s.FileID = string(data)
	return nil
}

type ExternalSource struct {
	URL string
}

func (ExternalSource) isFileSource() {}

func (s ExternalSource) MarsahlJSON() ([]byte, error) {
	return []byte(s.URL), nil
}

func (s *ExternalSource) UnmarshalJSON(data []byte) error {
	s.URL = string(data)
	return nil
}

type LocalSource struct {
	Filepath string
}

func (LocalSource) isFileSource() {}

const localSourcePrefix string = "attach://"

func (s LocalSource) MarsahlJSON() ([]byte, error) {
	return []byte(localSourcePrefix + s.Filepath), nil
}

func (s *LocalSource) UnmarshalJSON(data []byte) error {
	s.Filepath = strings.TrimPrefix(string(data), localSourcePrefix)
	return nil
}

type FileSource interface {
	isFileSource()
}
