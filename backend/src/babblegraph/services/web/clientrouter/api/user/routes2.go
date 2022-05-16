package user

import (
	"babblegraph/model/billing"
	"babblegraph/model/routes"
	"babblegraph/model/unsubscribereason"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/wordsmith"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// You should add new routes to this file

var Routes = router.RouteGroup{
	Prefix: "user",
	Routes: []router.Route{
		{
			Path: "signup_user_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.WithUTMEventTracking(
					"signup",
					routermiddleware.WithMaybePromotion(handleSignupUser),
				),
			),
		},
		{
			Path: "unsubscribe_user_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(unsubscribeUser),
			),
		}, {
			Path: "get_user_newsletter_preferences_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(getUserNewsletterPreferences),
			),
		}, {
			Path: "update_user_newsletter_preferences_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(updateUserNewsletterPreferences),
			),
		}, {
			Path: "upsert_user_vocabulary_entry_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(upsertUserVocabulary),
			),
		}, {
			Path: "get_user_vocabulary_entry_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(getUserVocabulary),
			),
		},
	},
}

type unsubscribeUserRequest struct {
	Token             string  `json:"token"`
	UnsubscribeReason *string `json:"unsubscribe_reason"`
	EmailAddress      *string `json:"email_address,omitempty"`
}

type unsubscribeUserResponse struct {
	Success bool                  `json:"success"`
	Error   *unsubscribeUserError `json:"error,omitempty"`
}

type unsubscribeUserError string

const (
	unsubscribeUserErrorMissingEmail   unsubscribeUserError = "missing-email"
	unsubscribeUserErrorIncorrectEmail unsubscribeUserError = "incorrect-email"
	unsubscribeUserErrorNoAuth         unsubscribeUserError = "no-auth"
	unsubscribeUserErrorIncorrectKey   unsubscribeUserError = "incorrect-key" // deliberately ambiguous
	unsubscribeUserErrorInvalidToken   unsubscribeUserError = "invalid-token"
)

func (u unsubscribeUserError) Ptr() *unsubscribeUserError {
	return &u
}

func unsubscribeUser(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req unsubscribeUserRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return unsubscribeUserResponse{
			Error: unsubscribeUserErrorInvalidToken.Ptr(),
		}, nil
	}
	if userAuth != nil {
		if *userID != userAuth.UserID {
			return unsubscribeUserResponse{
				Error: unsubscribeUserErrorIncorrectKey.Ptr(),
			}, nil
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := users.UnsubscribeUserByID(tx, *userID); err != nil {
				return err
			}
			if req.UnsubscribeReason != nil && len(*req.UnsubscribeReason) != 0 {
				err := unsubscribereason.InsertUnsubscribeReason(tx, *userID, wordsmith.LanguageCodeSpanish, *req.UnsubscribeReason)
				if err != nil {
					return err
				}
			}
			if userAuth.SubscriptionLevel != nil && *userAuth.SubscriptionLevel == useraccounts.SubscriptionLevelPremium {
				premiumNewsletterSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(r, tx, *userID)
				switch {
				case err != nil:
					return err
				case premiumNewsletterSubscription == nil:
					r.Infof("No active premium newsletter subscription for user %s", *userID)
				default:
					if err := billing.InsertPremiumNewsletterSyncRequest(tx, *premiumNewsletterSubscription.ID, billing.PremiumNewsletterSubscriptionUpdateTypeCanceled); err != nil {
						return err
					}
					return billing.CancelPremiumNewsletterSubscriptionForUser(r, tx, *userID)
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return unsubscribeUserResponse{
			Success: true,
		}, nil
	}
	if req.EmailAddress == nil {
		return unsubscribeUserResponse{
			Error: unsubscribeUserErrorMissingEmail.Ptr(),
		}, nil
	}
	formattedEmailAddress := email.FormatEmailAddress(*req.EmailAddress)
	var uErr *unsubscribeUserError
	err = database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.LookupUserForIDAndEmail(tx, *userID, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case user == nil:
			uErr = unsubscribeUserErrorIncorrectEmail.Ptr()
			return nil
		default:
			doesUserHaveAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, user.ID)
			switch {
			case err != nil:
				return err
			case doesUserHaveAccount:
				// The user has an account, but is not logged in
				uErr = unsubscribeUserErrorNoAuth.Ptr()
				return nil
			default:
				// no-op
			}
		}
		if req.UnsubscribeReason != nil && len(*req.UnsubscribeReason) != 0 {
			err := unsubscribereason.InsertUnsubscribeReason(tx, *userID, wordsmith.LanguageCodeSpanish, *req.UnsubscribeReason)
			if err != nil {
				return err
			}
		}
		premiumNewsletterSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(r, tx, *userID)
		switch {
		case err != nil:
			return err
		case premiumNewsletterSubscription == nil:
			r.Infof("No active premium newsletter subscription for user %s", *userID)
		default:
			if err := billing.InsertPremiumNewsletterSyncRequest(tx, *premiumNewsletterSubscription.ID, billing.PremiumNewsletterSubscriptionUpdateTypeCanceled); err != nil {
				return err
			}
			return billing.CancelPremiumNewsletterSubscriptionForUser(r, tx, *userID)
		}
		return users.UnsubscribeUserByID(tx, *userID)
	})
	switch {
	case err != nil:
		return nil, err
	case uErr != nil:
		if *uErr == unsubscribeUserErrorNoAuth {
			r.RespondWithStatus(http.StatusForbidden)
			return nil, nil
		}
		return unsubscribeUserResponse{
			Error: uErr,
		}, nil
	default:
		return unsubscribeUserResponse{
			Success: true,
		}, nil
	}
}
