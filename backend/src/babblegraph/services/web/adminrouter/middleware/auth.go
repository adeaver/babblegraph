package auth

import (
	"babblegraph/admin/model/auth"
	"babblegraph/admin/model/user"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type AuthenticatedRequestHandler func(adminID user.AdminID, r *router.Request) (interface{}, error)

func WithAuthentication(handler AuthenticatedRequestHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		for _, cookie := range r.GetCookies() {
			if cookie.Name == auth.AccessTokenCookieName {
				token := cookie.Value
				var adminID *user.AdminID
				if err := database.WithTx(func(tx *sqlx.Tx) error {
					var err error
					adminID, err = auth.ValidateAccessTokenAndGetUserID(tx, token)
					return err
				}); err != nil {
					r.RespondWithStatus(http.StatusForbidden)
					return nil, nil
				}
				return handler(*adminID, r)
			}
		}
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	}
}

func WithPermission(requiredPermission user.Permission, handler AuthenticatedRequestHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		for _, cookie := range r.GetCookies() {
			if cookie.Name == auth.AccessTokenCookieName {
				token := cookie.Value
				var adminID *user.AdminID
				if err := database.WithTx(func(tx *sqlx.Tx) error {
					var err error
					adminID, err = auth.ValidateAccessTokenAndGetUserID(tx, token)
					if err != nil {
						return err
					}
					return user.ValidateAdminUserPermission(tx, *adminID, requiredPermission)
				}); err != nil {
					r.RespondWithStatus(http.StatusForbidden)
					return nil, nil
				}
				return handler(*adminID, r)
			}
		}
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	}
}
