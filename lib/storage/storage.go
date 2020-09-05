package storage

import "fmt"

type EngineName string

const (
	EngineNameS3    EngineName = "S3"
	EngineNameLocal EngineName = "Local"
)

type engine interface {
	SaveDataToFile(filename *string, data []byte) (*Identifier, error)
	GetData(id Identifier) ([]byte, error)
}

func GetStorageEngineForName(name EngineName) engine {
	switch name {
	case EngineNameS3:
		panic("s3 storage is unimplemented")
	case EngineNameLocal:

	default:
		panic(fmt.Sprintf("unrecognized storage engine %s", name))
	}
}
