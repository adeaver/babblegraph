package auth

import (
	"babblegraph/model/admin"
	"babblegraph/model/routes"
	"babblegraph/services/web/adminrouter/middleware"
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
		}, {
			Path: "manage_user_permissions_1",
			Handler: middleware.WithPermission(
				admin.PermissionManagePermissions,
				manageUserPermissions,
			),
		}, {
			Path: "get_users_with_permissions_1",
			Handler: middleware.WithPermission(
				admin.PermissionManagePermissions,
				getUsersWithPermissions,
			),
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
		adminUser, err := admin.LookupAdminUserByEmailAddress(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case adminUser == nil:
			return nil
		}
		if err := admin.ValidateAdminUserPassword(tx, adminUser.ID, req.Password); err != nil {
			return nil
		}
		if err := admin.CreateTwoFactorAuthenticationAttempt(tx, adminUser.ID); err != nil {
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
		adminUser, err := admin.LookupAdminUserByEmailAddress(tx, formattedEmailAddress)
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
			if err := admin.ValidateTwoFactorAuthenticationAttempt(tx, adminUser.ID, req.TwoFactorAuthenticationCode); err != nil {
				return err
			}
		case env.EnvironmentLocal,
			env.EnvironmentLocalNoEmail,
			env.EnvironmentLocalTestEmail:
			// no-op
		default:
			return fmt.Errorf("Unrecognized environment: %s", envName)
		}
		accessToken, expirationTime, err = admin.CreateAccessToken(tx, adminUser.ID)
		return err
	}); err != nil {
		return nil, err
	}
	if accessToken == nil {
		return nil, fmt.Errorf("Null access token")
	}
	r.RespondWithCookie(&http.Cookie{
		Name:     admin.AccessTokenCookieName,
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
		if cookie.Name == admin.AccessTokenCookieName {
			token := cookie.Value
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return admin.InvalidateAccessToken(tx, token)
			}); err != nil {
				return nil, err
			}
			r.RespondWithCookie(&http.Cookie{
				Name:     admin.AccessTokenCookieName,
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
	var adminUserID admin.ID
	if err := encrypt.WithDecodedToken(req.Token, func(t encrypt.TokenPair) error {
		if t.Key != routes.AdminRegistrationKey.Str() {
			return fmt.Errorf("incorrect key")
		}
		adminIDStr, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("incorrect type")
		}
		adminUserID = admin.ID(adminIDStr)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		adminUser, err := admin.GetAdminUser(tx, adminUserID)
		switch {
		case err != nil:
			return err
		case adminUser.EmailAddress != formattedEmailAddress:
			return fmt.Errorf("Invalid")
		}
		if err := admin.CreateAdminUserPassword(tx, adminUser.ID, req.Password); err != nil {
			return err
		}
		if err := admin.CreateTwoFactorAuthenticationAttempt(tx, adminUser.ID); err != nil {
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
	var adminUserID admin.ID
	if err := encrypt.WithDecodedToken(req.Token, func(t encrypt.TokenPair) error {
		if t.Key != routes.AdminRegistrationKey.Str() {
			return fmt.Errorf("incorrect key")
		}
		adminIDStr, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("incorrect type")
		}
		adminUserID = admin.ID(adminIDStr)
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
			if err := admin.ValidateTwoFactorAuthenticationAttempt(tx, adminUserID, req.TwoFactorAuthenticationCode); err != nil {
				return err
			}
		case env.EnvironmentLocal,
			env.EnvironmentLocalNoEmail,
			env.EnvironmentLocalTestEmail:
			// no-op
		default:
			return fmt.Errorf("Unrecognized environment: %s", envName)
		}
		if err := admin.ActivateAdminUserPassword(tx, adminUserID); err != nil {
			return err
		}
		var err error
		accessToken, expirationTime, err = admin.CreateAccessToken(tx, adminUserID)
		return err
	}); err != nil {
		return nil, err
	}
	if accessToken == nil {
		return nil, fmt.Errorf("No access token")
	}
	r.RespondWithCookie(&http.Cookie{
		Name:     admin.AccessTokenCookieName,
		Value:    *accessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  *expirationTime,
	})
	return validateTwoFactorAuthenticationCodeForCreateResponse{
		Success: true,
	}, nil
}

type manageUserPermissionsRequest struct {
	AdminID admin.ID           `json:"admin_id"`
	Updates []permissionUpdate `json:"updates"`
}

type permissionUpdate struct {
	Permission admin.Permission `json:"permission"`
	IsActive   bool             `json:"is_active"`
}

type manageUserPermissionResponse struct {
	Success bool `json:"success"`
}

func manageUserPermissions(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req manageUserPermissionsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		for _, u := range req.Updates {
			if err := admin.UpsertUserPermission(tx, req.AdminID, u.Permission, u.IsActive); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return manageUserPermissionResponse{
		Success: true,
	}, nil
}

type getUsersWithPermissionsResponse struct {
	Users []userWithPermissions `json:"users"`
}

type userWithPermissions struct {
	ID           admin.ID           `json:"id"`
	EmailAddress string             `json:"email_address"`
	Permissions  []admin.Permission `json:"permissions"`
}

func getUsersWithPermissions(adminID admin.ID, r *router.Request) (interface{}, error) {
	var permissionMappings []admin.UserPermissionMapping
	var adminUsers []admin.Admin
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		adminUsers, err = admin.GetAllAdminUsers(tx)
		if err != nil {
			return err
		}
		permissionMappings, err = admin.GetAllActiveUserPermissions(tx)
		return err
	}); err != nil {
		return nil, err
	}
	adminPermissionsToID := make(map[admin.ID][]admin.Permission)
	for _, p := range permissionMappings {
		adminPermissionsToID[p.AdminUserID] = append(adminPermissionsToID[p.AdminUserID], p.Permission)
	}
	var mappings []userWithPermissions
	for _, u := range adminUsers {
		mappings = append(mappings, userWithPermissions{
			ID:           u.ID,
			EmailAddress: u.EmailAddress,
			Permissions:  adminPermissionsToID[u.ID],
		})
	}
	return getUsersWithPermissionsResponse{
		Users: mappings,
	}, nil
}
