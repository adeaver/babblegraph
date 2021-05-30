package random

import (
	"babblegraph/util/ptr"
	"fmt"
	"math/rand"
	"time"
)

const (
	stringCharSet = "abcdefghijklmnopqrstuvwxyz"
)

func MakeRandomString(inputLength int) (*string, error) {
	if inputLength <= 0 {
		return nil, fmt.Errorf("Must have at least length of 1")
	}
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var chars []byte
	for i := 0; i < inputLength; i++ {
		b := stringCharSet[seededRand.Intn(len(stringCharSet))]
		if seededRand.Intn(100) > 50 {
			b -= 32
		}
		chars = append(chars, b)
	}
	return ptr.String(string(chars)), nil
}

func MustMakeRandomString(inputLength int) string {
	s, err := MakeRandomString(inputLength)
	if err != nil {
		panic(err.Error())
	}
	return *s
}
