package virtualfile

import (
	"babblegraph/util/encrypt"
	"babblegraph/util/ptr"
	"fmt"
	"strings"
)

type Type string

const (
	TypePodcast      Type = "podcast"
	TypePodcastImage Type = "podcast-image"
)

func (t Type) Str() string {
	return string(t)
}

func (t Type) Ptr() *Type {
	return &t
}

func typeFromString(t string) (*Type, error) {
	switch strings.ToLower(t) {
	case strings.ToLower(TypePodcast.Str()):
		return TypePodcast.Ptr(), nil
	case strings.ToLower(TypePodcastImage.Str()):
		return TypePodcastImage.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unrecognized virutal file type %s", t)
	}
}

func EncodeAsVirtualFileWithType(id string, t Type) (*string, error) {
	return encrypt.GetToken(encrypt.TokenPair{
		Key:   t.Str(),
		Value: id,
	})
}

func GetObjectIDAndType(virtualFileKey string) (*string, *Type, error) {
	var objectID *string
	var t *Type
	if err := encrypt.WithDecodedToken(virtualFileKey, func(token encrypt.TokenPair) error {
		var err error
		t, err = typeFromString(token.Key)
		if err != nil {
			return err
		}
		id, ok := token.Value.(string)
		if !ok {
			return fmt.Errorf("URL did not decode correctly")
		}
		objectID = ptr.String(id)
		return nil
	}); err != nil {
		return nil, nil, err
	}
	return objectID, t, nil
}
