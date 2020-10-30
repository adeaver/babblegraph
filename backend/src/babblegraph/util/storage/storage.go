package storage

import (
	"babblegraph/util/env"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

type storageType string

const (
	storageTypeLocal storageType = "local"
)

type FileIdentifier string

func (f FileIdentifier) Str() string {
	return string(f)
}

func (f FileIdentifier) Ptr() *FileIdentifier {
	return &f
}

func (f FileIdentifier) ToFileName() string {
	folder := env.GetEnvironmentVariableOrDefault("STORAGE_LOCATION", "/tmp")
	return fmt.Sprintf("%s/%s", folder, f.Str())
}

func ReadFile(id FileIdentifier) ([]byte, error) {
	return ioutil.ReadFile(id.ToFileName())
}

func WriteFile(fileType string, body string) (*FileIdentifier, error) {
	id := FileIdentifier(fmt.Sprintf("%s.%s", uuid.New(), fileType))
	if err := ioutil.WriteFile(id.ToFileName(), []byte(body), 0644); err != nil {
		return nil, err
	}
	return id.Ptr(), nil
}

func DeleteFile(id FileIdentifier) error {
	return os.Remove(id.ToFileName())
}
