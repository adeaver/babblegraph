package auth

import (
	"babblegraph/admin/model/auth"
	"babblegraph/admin/model/user"
	"babblegraph/model/routes"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/util/encrypt"
	"babblegraph/util/env"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "auth",
	Routes: []router.Route{
		{
			Path:    "validate_login_credentials_1",
			Handler: validateLoginCredentials,
		}, {
			Path:    "validate_two_factor_code_1",
			Handler: validateTwoFactorAuthenticationCode,
		}, {
			Path:    "invalidate_login_credentials_1",
			Handler: invalidateCredentials,
		}, {
			Path:    "create_admin_user_password_1",
			Handler: createAdminUserPassword,
		}, {
			Path:    "validate_two_factor_code_for_create_1",
			Handler: validateTwoFactorAuthenticationCodeForCreate,
		},
	},
}

type validateLoginCredentialsRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

type validateLoginCredentialsResponse struct {
	Success bool `json:"success"`
}

func validateLoginCredentials(r *router.Request) (interface{}, error) {
	var req validateLoginCredentialsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var success bool
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		adminUser, err := user.LookupAdminUserByEmailAddress(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case adminUser == nil:
			return nil
		}
		if err := user.ValidateAdminUserPassword(tx, adminUser.AdminID, req.Password); err != nil {
			return nil
		}
		if err := auth.CreateTwoFactorAuthenticationAttempt(tx, adminUser.AdminID); err != nil {
			return err
		}
		success = true
		return nil
	}); err != nil {
		return nil, err
	}
	return validateLoginCredentialsResponse{
		Success: success,
	}, nil
}

type validateTwoFactorAuthenticationCodeRequest struct {
	EmailAddress                string `json:"email_address"`
	TwoFactorAuthenticationCode string `json:"two_factor_authentication_code"`
}

type validateTwoFactorAuthenticationCodeResponse struct {
	Success bool `json:"success"`
}

func validateTwoFactorAuthenticationCode(r *router.Request) (interface{}, error) {
	var req validateTwoFactorAuthenticationCodeRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	var accessToken *string
	var expirationTime *time.Time
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		adminUser, err := user.LookupAdminUserByEmailAddress(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case adminUser == nil:
			return nil
		}
		envName := env.MustEnvironmentName()
		switch envName {
		case env.EnvironmentProd,
			env.EnvironmentStage:
			if err := auth.ValidateTwoFactorAuthenticationAttempt(tx, adminUser.AdminID, req.TwoFactorAuthenticationCode); err != nil {
				return err
			}
		case env.EnvironmentLocal,
			env.EnvironmentLocalNoEmail,
			env.EnvironmentLocalTestEmail:
			// no-op
		default:
			return fmt.Errorf("Unrecognized environment: %s", envName)
		}
		accessToken, expirationTime, err = auth.CreateAccessToken(tx, adminUser.AdminID)
		return err
	}); err != nil {
		return nil, err
	}
	if accessToken == nil {
		return nil, fmt.Errorf("Null access token")
	}
	r.RespondWithCookie(&http.Cookie{
		Name:     auth.AccessTokenCookieName,
		Value:    *accessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  *expirationTime,
	})
	return validateTwoFactorAuthenticationCodeResponse{
		Success: true,
	}, nil
}

type invalidateCredentialsRequest struct{}

type invalidateCredentialsResponse struct {
	Success bool `json:"success"`
}

func invalidateCredentials(r *router.Request) (interface{}, error) {
	for _, cookie := range r.GetCookies() {
		if cookie.Name == auth.AccessTokenCookieName {
			token := cookie.Value
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return auth.InvalidateAccessToken(tx, token)
			}); err != nil {
				return nil, err
			}
			r.RespondWithCookie(&http.Cookie{
				Name:     auth.AccessTokenCookieName,
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(-5 * time.Hour),
			})
		}
	}
	return invalidateCredentialsResponse{
		Success: true,
	}, nil
}

type createAdminUserPasswordRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
	Token        string `json:"token"`
}

type createAdminUserPasswordResponse struct {
	Success bool `json:"success"`
}

func createAdminUserPassword(r *router.Request) (interface{}, error) {
	var req createAdminUserPasswordRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	var adminUserID user.AdminID
	if err := encrypt.WithDecodedToken(req.Token, func(t encrypt.TokenPair) error {
		if t.Key != routes.AdminRegistrationKey.Str() {
			return fmt.Errorf("incorrect key")
		}
		adminIDStr, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("incorrect type")
		}
		adminUserID = user.AdminID(adminIDStr)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		adminUser, err := user.GetAdminUser(tx, adminUserID)
		switch {
		case err != nil:
			return err
		case adminUser.EmailAddress != formattedEmailAddress:
			return fmt.Errorf("Invalid")
		}
		if err := user.CreateAdminUserPassword(tx, adminUser.AdminID, req.Password); err != nil {
			return err
		}
		if err := auth.CreateTwoFactorAuthenticationAttempt(tx, adminUser.AdminID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return createAdminUserPasswordResponse{
		Success: true,
	}, nil
}

type validateTwoFactorAuthenticationCodeForCreateRequest struct {
	Token                       string `json:"token"`
	TwoFactorAuthenticationCode string `json:"two_factor_authentication_code"`
}

type validateTwoFactorAuthenticationCodeForCreateResponse struct {
	Success bool `json:"success"`
}

func validateTwoFactorAuthenticationCodeForCreate(r *router.Request) (interface{}, error) {
	var req validateTwoFactorAuthenticationCodeForCreateRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var adminUserID user.AdminID
	if err := encrypt.WithDecodedToken(req.Token, func(t encrypt.TokenPair) error {
		if t.Key != routes.AdminRegistrationKey.Str() {
			return fmt.Errorf("incorrect key")
		}
		adminIDStr, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("incorrect type")
		}
		adminUserID = user.AdminID(adminIDStr)
		return nil
	}); err != nil {
		return nil, err
	}
	var accessToken *string
	var expirationTime *time.Time
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		envName := env.MustEnvironmentName()
		switch envName {
		case env.EnvironmentProd,
			env.EnvironmentStage:
			if err := auth.ValidateTwoFactorAuthenticationAttempt(tx, adminUserID, req.TwoFactorAuthenticationCode); err != nil {
				return err
			}
		case env.EnvironmentLocal,
			env.EnvironmentLocalNoEmail,
			env.EnvironmentLocalTestEmail:
			// no-op
		default:
			return fmt.Errorf("Unrecognized environment: %s", envName)
		}
		if err := user.ActivateAdminUserPassword(tx, adminUserID); err != nil {
			return err
		}
		var err error
		accessToken, expirationTime, err = auth.CreateAccessToken(tx, adminUserID)
		return err
	}); err != nil {
		return nil, err
	}
	if accessToken == nil {
		return nil, fmt.Errorf("No access token")
	}
	r.RespondWithCookie(&http.Cookie{
		Name:     auth.AccessTokenCookieName,
		Value:    *accessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  *expirationTime,
	})
	return validateTwoFactorAuthenticationCodeForCreateResponse{
		Success: true,
	}, nil
}
