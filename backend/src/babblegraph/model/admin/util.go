package user

import (
	"babblegraph/util/ptr"
	"babblegraph/util/random"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	defaultSaltLength                        = 32
	defaultTwoFactorAuthenticationCodeLength = 8
	defaultAccessTokenLength                 = 64
)

func generatePasswordSalt() (*string, error) {
	saltBytes := make([]byte, defaultSaltLength)
	if _, err := rand.Read(saltBytes); err != nil {
		return nil, err
	}
	return ptr.String(base64.URLEncoding.EncodeToString(saltBytes)), nil
}

func generatePasswordHash(password, salt string) (*string, error) {
	saltedPassword := makeSaltedPassword(password, salt)
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

// Valid passwords contain:
// Between 16 and 32 characters
// At least 3 of: Capital Letter, Lowercase Letter, Number, Symbol
func validatePasswordMeetsRequirements(password string) bool {
	if len(password) < 16 || len(password) > 32 {
		return false
	}
	requirements := []int{0, 0, 0, 0}
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			requirements[0] = 1
		case unicode.IsLower(c):
			requirements[1] = 1
		case unicode.IsNumber(c):
			requirements[2] = 1
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			requirements[3] = 1
		default:
			return false
		}
	}
	total := 0
	for _, i := range requirements {
		total += i
	}
	return total >= 3
}

func generateTwoFactorAuthenticationCode() string {
	return random.MustMakeRandomString(defaultTwoFactorAuthenticationCodeLength)
}

func generateAccessToken() string {
	return random.MustMakeRandomString(defaultAccessTokenLength)
}
