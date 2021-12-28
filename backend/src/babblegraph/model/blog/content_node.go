package blog

import "encoding/json"

type Content []ContentNode

type ContentNode interface {
	json.Marshaler
	json.Unmarshaler
	GetType() ContentNodeType
}

type ContentNodeType string
