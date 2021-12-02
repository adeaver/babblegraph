package middleware

import (
	"babblegraph/model/admin"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type AuthenticatedRequestHandler func(adminID admin.ID, r *router.Request) (interface{}, error)

func WithAuthentication(handler AuthenticatedRequestHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		for _, cookie := range r.GetCookies() {
			if cookie.Name == admin.AccessTokenCookieName {
				token := cookie.Value
				var adminID *admin.ID
				if err := database.WithTx(func(tx *sqlx.Tx) error {
					var err error
					adminID, err = admin.ValidateAccessTokenAndGetUserID(tx, token)
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

func WithPermission(requiredPermission admin.Permission, handler AuthenticatedRequestHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		for _, cookie := range r.GetCookies() {
			if cookie.Name == admin.AccessTokenCookieName {
				token := cookie.Value
				var adminID *admin.ID
				if err := database.WithTx(func(tx *sqlx.Tx) error {
					var err error
					adminID, err = admin.ValidateAccessTokenAndGetUserID(tx, token)
					if err != nil {
						return err
					}
					return admin.ValidateAdminUserPermission(tx, *adminID, requiredPermission)
				}); err != nil {
					log.Println(err.Error())
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
