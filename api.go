package telegrambot

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

type InternalSource struct {
	FileID string
}

func (InternalSource) isFileSource() {}

func (s InternalSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.FileID)
}

func (s *InternalSource) UnmarshalJSON(data []byte) error {
	s.FileID = string(data)
	return nil
}

type ExternalSource struct {
	URL string
}

func (ExternalSource) isFileSource() {}

func (s ExternalSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.URL)
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

func filenameFromFilepath(path string) string {
	filename := filepath.Base(path)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func (s LocalSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(localSourcePrefix + filenameFromFilepath(s.Filepath))
}

func (s *LocalSource) UnmarshalJSON(data []byte) error {
	s.Filepath = strings.TrimPrefix(string(data), localSourcePrefix)
	return nil
}

type FileSource interface {
	isFileSource()
}
