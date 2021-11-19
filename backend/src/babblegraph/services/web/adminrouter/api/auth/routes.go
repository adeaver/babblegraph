package auth

import (
	"babblegraph/admin/model/auth"
	"babblegraph/admin/model/user"
	"babblegraph/model/routes"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/util/encrypt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "auth",
	Routes: []router.Route{
		{
			Path: "validate_login_credentials_1",
			Handler: router.RouteHandler{
				HandleRequestBody: validateLoginCredentials,
			},
		}, {
			Path: "validate_two_factor_code_1",
			Handler: router.RouteHandler{
				HandleRawRequest: validateTwoFactorAuthenticationCode,
			},
		}, {
			Path: "invalidate_login_credentials_1",
			Handler: router.RouteHandler{
				HandleRawRequest: invalidateCredentials,
			},
		}, {
			Path: "create_admin_user_password_1",
			Handler: router.RouteHandler{
				HandleRequestBody: createAdminUserPassword,
			},
		}, {
			Path: "validate_two_factor_code_for_create_1",
			Handler: router.RouteHandler{
				HandleRawRequest: validateTwoFactorAuthenticationCodeForCreate,
			},
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

func validateLoginCredentials(reqBody []byte) (interface{}, error) {
	var req validateLoginCredentialsRequest
	if err := json.Unmarshal(reqBody, &req); err != nil {
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
		// TODO: validate password
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

func validateTwoFactorAuthenticationCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req validateTwoFactorAuthenticationCodeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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
		if err := auth.ValidateTwoFactorAuthenticationAttempt(tx, adminUser.AdminID, req.TwoFactorAuthenticationCode); err != nil {
			return err
		}
		accessToken, expirationTime, err = auth.CreateAccessToken(tx, adminUser.AdminID)
		return err
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if accessToken == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     auth.AccessTokenCookieName,
		Value:    *accessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  *expirationTime,
	})
	json.NewEncoder(w).Encode(validateTwoFactorAuthenticationCodeResponse{
		Success: true,
	})
}

type invalidateCredentialsRequest struct{}

type invalidateCredentialsResponse struct {
	Success bool `json:"success"`
}

func invalidateCredentials(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	for _, cookie := range r.Cookies() {
		if cookie.Name == auth.AccessTokenCookieName {
			token := cookie.Value
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return auth.InvalidateAccessToken(tx, token)
			}); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     auth.AccessTokenCookieName,
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(-5 * time.Hour),
			})
		}
	}
	json.NewEncoder(w).Encode(invalidateCredentialsResponse{
		Success: true,
	})
}

type createAdminUserPasswordRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
	Token        string `json:"token"`
}

type createAdminUserPasswordResponse struct {
	Success bool `json:"success"`
}

func createAdminUserPassword(body []byte) (interface{}, error) {
	var req createAdminUserPasswordRequest
	if err := json.Unmarshal(body, &req); err != nil {
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

func validateTwoFactorAuthenticationCodeForCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req validateTwoFactorAuthenticationCodeForCreateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var accessToken *string
	var expirationTime *time.Time
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		if err := auth.ValidateTwoFactorAuthenticationAttempt(tx, adminUserID, req.TwoFactorAuthenticationCode); err != nil {
			return err
		}
		if err := user.ActivateAdminUserPassword(tx, adminUserID); err != nil {
			return err
		}
		accessToken, expirationTime, err = auth.CreateAccessToken(tx, adminUserID)
		return err
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if accessToken == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     auth.AccessTokenCookieName,
		Value:    *accessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  *expirationTime,
	})
	json.NewEncoder(w).Encode(validateTwoFactorAuthenticationCodeForCreateResponse{
		Success: true,
	})
}
