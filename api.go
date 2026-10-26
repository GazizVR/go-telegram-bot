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
	var filepath string
	if err := json.Unmarshal(data, &filepath); err != nil {
		return err
	}
	s.Filepath = strings.TrimPrefix(filepath, localSourcePrefix)
	return nil
}

type FileSource interface {
	isFileSource()
}
