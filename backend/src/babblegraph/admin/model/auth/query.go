package auth

import (
	"babblegraph/admin/model/user"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	create2FACodeQuery        = "INSERT INTO admin_2fa_codes (code, admin_user_id) VALUES ($1, $2)"
	expireNonActiveCodesQuery = "UPDATE admin_2fa_codes SET expires_at = timezone('utc', now()) WHERE admin_user_id = $1"
	get2FACodeQuery           = "SELECT * FROM admin_2fa_codes WHERE admin_user_id = $1 AND code = $2"

	getUnfulfilled2FACodeQuery   = "SELECT * FROM admin_2fa_codes WHERE expires_at IS NULL"
	update2FACodeExpirationQuery = "UPDATE admin_2fa_codes SET expires_at = $1 WHERE admin_user_id = $2 AND code = $3 AND expires_at IS NULL"
	deleteExpired2FACodesQuery   = "DELETE FROM admin_2fa_codes WHERE expires_at < current_date - interval '2 days'"

	default2FAExpirationTime = 10 * time.Minute

	createAccessTokenQuery = `INSERT INTO
        admin_access_token (
            token,
            expires_at,
            admin_user_id
        ) VALUES ($1, $2, $3)
        ON CONFLICT (admin_user_id) DO UPDATE
        SET expires_at = $2, token = $1`
	getAccessTokenQuery            = "SELECT * FROM admin_access_token WHERE token = $1"
	deleteAccessTokenQuery         = "DELETE FROM admin_access_token WHERE token = $1"
	deleteExpiredAccessTokensQuery = "DELETE FROM admin_access_token WHERE expires_at < current_date - interval '2 days'"

	defaultAccessTokenExpirationTime = 24 * time.Hour
)

// Public functions

func CreateTwoFactorAuthenticationAttempt(tx *sqlx.Tx, adminUserID user.AdminID) error {
	code := generateTwoFactorAuthenticationCode()
	if _, err := tx.Exec(expireNonActiveCodesQuery, adminUserID); err != nil {
		return err
	}
	if _, err := tx.Exec(create2FACodeQuery, code, adminUserID); err != nil {
		return err
	}
	return nil
}

func ValidateTwoFactorAuthenticationAttempt(tx *sqlx.Tx, adminUserID user.AdminID, code string) error {
	var codes []dbAdmin2FACode
	err := tx.Select(&codes, get2FACodeQuery, adminUserID, code)
	switch {
	case err != nil:
		return err
	case len(codes) == 0:
		return fmt.Errorf("Invalid code")
	default:
		for _, c := range codes {
			if c.ExpiresAt != nil && time.Now().Before(*c.ExpiresAt) {
				return nil
			}
		}
	}
	return fmt.Errorf("Code expired")
}

func GetUnfulfilledTwoFactorAuthenticationAttempts(tx *sqlx.Tx) ([]Admin2FACode, error) {
	var codes []dbAdmin2FACode
	if err := tx.Select(&codes, getUnfulfilled2FACodeQuery); err != nil {
		return nil, err
	}
	var out []Admin2FACode
	for _, c := range codes {
		out = append(out, c.ToNonDB())
	}
	return out, nil
}

func FulfillTwoFactorAuthenticationAttempt(tx *sqlx.Tx, adminID user.AdminID, code string) error {
	if _, err := tx.Exec(update2FACodeExpirationQuery, time.Now().Add(default2FAExpirationTime), adminID, code); err != nil {
		return err
	}
	return nil
}

func RemoveExpiredTwoFactorCodes(tx *sqlx.Tx) error {
	if _, err := tx.Exec(deleteExpired2FACodesQuery); err != nil {
		return err
	}
	return nil
}

func CreateAccessToken(tx *sqlx.Tx, adminID user.AdminID) (*string, *time.Time, error) {
	accessToken := generateAccessToken()
	expirationTime := time.Now().Add(defaultAccessTokenExpirationTime)
	if _, err := tx.Exec(createAccessTokenQuery, accessToken, expirationTime, adminID); err != nil {
		return nil, nil, err
	}
	return &accessToken, &expirationTime, nil
}

func ValidateAccessTokenAndGetUserID(tx *sqlx.Tx, accessToken string) (*user.AdminID, error) {
	var accessTokens []dbAdminAccessToken
	err := tx.Select(&accessTokens, getAccessTokenQuery, accessToken)
	switch {
	case err != nil:
		return nil, err
	case len(accessTokens) == 0:
		return nil, fmt.Errorf("Invalid access token")
	case len(accessTokens) > 1:
		return nil, fmt.Errorf("Cannot disambiguate access tokens")
	default:
		if accessTokens[0].ExpiresAt.Before(time.Now()) {
			return nil, fmt.Errorf("Expired access token")
		}
		return &accessTokens[0].AdminUserID, nil
	}
}

func InvalidateAccessToken(tx *sqlx.Tx, accessToken string) error {
	if _, err := tx.Exec(deleteAccessTokenQuery, accessToken); err != nil {
		return err
	}
	return nil
}

func RemoveExpiredAccessTokens(tx *sqlx.Tx) error {
	if _, err := tx.Exec(deleteExpiredAccessTokensQuery); err != nil {
		return err
	}
	return nil
}
