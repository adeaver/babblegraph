package auth

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const expirationTime3days = 3 * 24 * time.Hour

func CreateJWTForUser(userID users.UserID) (*string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expirationTime3days).Unix(),
		Subject:   string(userID),
		IssuedAt:  time.Now().Unix(),
	})
	signedToken, err := accessToken.SignedString(env.MustEnvironmentVariable("HMAC_SECRET"))
	if err != nil {
		return nil, err
	}
	return ptr.String(signedToken), nil
}

func VerifyJWTAndGetUserID(tokenString string) (*users.UserID, bool, error) {
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return env.MustEnvironmentVariable("HMAC_SECRET"), nil
	})
	if err != nil {
		return nil, false, err
	}
	if !token.Valid {
		return nil, false, nil
	}
	userID := users.UserID(token.Subject)
	return &userID, true, nil
}
