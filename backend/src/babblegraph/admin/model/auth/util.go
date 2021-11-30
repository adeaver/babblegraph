package auth

import "babblegraph/util/random"

const (
	defaultTwoFactorAuthenticationCodeLength = 8
	defaultAccessTokenLength                 = 64
)

func generateTwoFactorAuthenticationCode() string {
	return random.MustMakeRandomString(defaultTwoFactorAuthenticationCodeLength)
}

func generateAccessToken() string {
	return random.MustMakeRandomString(defaultAccessTokenLength)
}
