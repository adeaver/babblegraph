package encrypt

import (
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
)

// A token pair is used so different
// services can use the same data but not encrypt
// to the same value.

// i.e. Unsubscribing can be done via /unsubscribe/:token
// where token is generated with the key "unsubscribe"
// AND changing preferences can be done with the key "preferences"
type TokenPair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func GetToken(t TokenPair) (*string, error) {
	encryptionKey := env.MustEnvironmentVariable("AES_KEY")
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, err
	}
	jsonToken, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(jsonToken))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], jsonToken)
	return ptr.String(base64.URLEncoding.EncodeToString(ciphertext)), nil
}

func WithDecodedToken(token string, fn func(TokenPair) error) error {
	log.Println(fmt.Sprintf("Decoding token %s", token))
	encryptionKey := env.MustEnvironmentVariable("AES_KEY")
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return err
	}
	if len(token) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}
	decodedToken, err := base64.URLEncoding.DecodeString(string(token))
	if err != nil {
		return err
	}
	iv := decodedToken[:aes.BlockSize]
	decodedToken = decodedToken[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decodedToken, decodedToken)
	var t TokenPair
	if err := json.Unmarshal(decodedToken, &t); err != nil {
		return err
	}
	return fn(t)
}
