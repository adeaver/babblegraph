package virtualfile

import (
	"babblegraph/util/encrypt"
	"babblegraph/util/ptr"
	"fmt"
	"strings"
)

type Type string

const (
	TypePodcast Type = "podcast"
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
	default:
		return nil, fmt.Errorf("Unrecognized virutal file type %s", t)
	}
}

func EncodeAsVirtualFileWithType(url string, t Type) (*string, error) {
	return encrypt.GetToken(encrypt.TokenPair{
		Key:   t.Str(),
		Value: url,
	})
}

func GetURLAndType(virtualFileKey string) (*string, *Type, error) {
	var url *string
	var t *Type
	if err := encrypt.WithDecodedToken(virtualFileKey, func(token encrypt.TokenPair) error {
		var err error
		t, err = typeFromString(token.Key)
		if err != nil {
			return err
		}
		u, ok := token.Value.(string)
		if !ok {
			return fmt.Errorf("URL did not decode correctly")
		}
		url = ptr.String(u)
		return nil
	}); err != nil {
		return nil, nil, err
	}
	return url, t, nil
}
