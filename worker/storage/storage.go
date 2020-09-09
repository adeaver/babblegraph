package storage

import (
	"fmt"
	"io/ioutil"

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
	return fmt.Sprintf("/tmp/%s", f.Str())
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
