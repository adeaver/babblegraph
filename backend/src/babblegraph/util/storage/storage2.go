package storage

import (
	"fmt"
	"strings"
)

func init() {
	RemoteStorage = getRemoteStorageForEnvironment()
}

type Storage interface {
	Read(directory string, fileName string) ([]byte, error)
	Write(directory string, file File) error
	Delete(directory string, fileName string) error
	DoesExist(directory string, fileName string) (bool, error)
}

type File struct {
	name          string
	data          []byte
	contentType   ContentType
	accessControl AccessControl
}

func NewFile(fileName string, data []byte) (*File, error) {
	contentType, err := getContentTypeForFileName(fileName)
	if err != nil {
		return nil, err
	}
	return &File{
		name:          fileName,
		data:          data,
		contentType:   *contentType,
		accessControl: AccessControlPrivate,
	}, nil
}

func (f *File) AssignAccessControlLevel(accessControl AccessControl) {
	f.accessControl = accessControl
}

type ContentType string

const (
	ContentTypeApplicationJSON ContentType = "application/json"

	ContentTypeImageJPEG ContentType = "image/jpeg"
)

func (c ContentType) Str() string {
	return string(c)
}

func (c ContentType) Ptr() *ContentType {
	return &c
}

func getContentTypeForFileName(fileName string) (*ContentType, error) {
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) <= 1 {
		return nil, fmt.Errorf("Invalid file name")
	}
	extension := strings.Join(fileNameParts[1:], ".")
	switch {
	case extension == "json":
		return ContentTypeApplicationJSON.Ptr(), nil
	case extension == "jpeg":
		return ContentTypeImageJPEG.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unrecognized file type %s", extension)
	}
}

type AccessControl string

const (
	AccessControlPrivate        AccessControl = "private"
	AccessControlPublicReadOnly AccessControl = "public-read"
)

func (a AccessControl) Str() string {
	return string(a)
}
