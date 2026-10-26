package telegrambot

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

type RemoteSource struct {
	FileIDOrURL string
}

func (RemoteSource) isFileSource() {}

func (s RemoteSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.FileIDOrURL)
}

func (s *RemoteSource) UnmarshalJSON(data []byte) error {
	var fileIDOrURL string
	if err := json.Unmarshal(data, &fileIDOrURL); err != nil {
		return err
	}
	s.FileIDOrURL = fileIDOrURL
	return nil
}

type sourcePrefix string

const (
	attachPrefix sourcePrefix = "attach://"
	filePrefix   sourcePrefix = "file://"
)

type LocalSource struct {
	prefix   sourcePrefix
	Filepath string
}

func (*LocalSource) isFileSource() {}

func filenameFromFilepath(path string) string {
	filename := filepath.Base(path)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func (s LocalSource) MarshalJSON() ([]byte, error) {
	filepath := s.Filepath
	if s.prefix == attachPrefix {
		filepath = filenameFromFilepath(s.Filepath)
	}
	return json.Marshal(string(s.prefix) + filepath)
}

func (s *LocalSource) UnmarshalJSON(data []byte) error {
	var rawFilepath string
	if err := json.Unmarshal(data, &rawFilepath); err != nil {
		return err
	}
	prefixInd := strings.Index(rawFilepath, "://")
	s.prefix = sourcePrefix(rawFilepath[:prefixInd])
	filepathInd := prefixInd + 3
	s.Filepath = rawFilepath[filepathInd:]
	return nil
}

type FileSource interface {
	isFileSource()
}
