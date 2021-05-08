package useraccounts

import (
	"babblegraph/util/ptr"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const defaultSaltLength = 16

func generatePasswordSalt() (*string, error) {
	saltBytes := make([]byte, defaultSaltLength)
	if _, err := rand.Read(saltBytes); err != nil {
		return nil, err
	}
	return ptr.String(base64.URLEncoding.EncodeToString(saltBytes)), nil
}

func generatePasswordHash(password, salt string) (*string, error) {
	saltedPassword := makeSaltedPassword(password, salt)
	log.Println(saltedPassword)
	hash, err := bcrypt.GenerateFromPassword(saltedPassword, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return ptr.String(string(hash)), nil
}

func makeSaltedPassword(password, salt string) []byte {
	return []byte(fmt.Sprintf("%s%s", password, salt))
}

func comparePasswords(hashedPassword, password, salt string) error {
	saltedPassword := makeSaltedPassword(password, salt)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), saltedPassword)
}
