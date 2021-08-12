package router

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/useraccountsnotifications"
	"babblegraph/model/users"
	"babblegraph/services/web/middleware"
	"babblegraph/services/web/util/routetoken"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/util/encrypt"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func registerUserAccountsRoutes() {
	a.prefixes["useraccounts"] = true

	a.r.HandleFunc("/api/useraccounts/login_user_1", middleware.WithoutBodyLogger(loginUser))
	a.routeNames["/api/useraccounts/login_user_1"] = true

	a.r.HandleFunc("/api/useraccounts/create_user_1", middleware.WithoutBodyLogger(createUser))
	a.routeNames["/api/useraccounts/create_user_1"] = true

	a.r.HandleFunc("/api/useraccounts/reset_password_1", middleware.WithoutBodyLogger(resetPassword))
	a.routeNames["/api/useraccounts/reset_password_1"] = true

	a.r.HandleFunc("/api/useraccounts/get_user_profile_1", middleware.WithoutBodyLogger(getUserProfile))
	a.routeNames["/api/useraccounts/get_user_profile_1"] = true
}

type loginUserRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
	RedirectKey  string `json:"redirect_key"`
}

type loginUserResponse struct {
	Location   *string     `json:"location,omitempty"`
	LoginError *loginError `json:"login_error"`
}

type loginError string

const (
	loginErrorInvalidCredentials loginError = "invalid-creds"
)

func (l loginError) Ptr() *loginError {
	return &l
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	var req loginUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	var userID *users.UserID
	var lErr *loginError
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case user == nil:
			lErr = loginErrorInvalidCredentials.Ptr()
			return fmt.Errorf("no user found for email address")
		}
		userID = &user.ID
		err = useraccounts.VerifyPasswordForUser(tx, user.ID, req.Password)
		if err != nil {
			lErr = loginErrorInvalidCredentials.Ptr()
			return err
		}
		return nil
	}); err != nil {
		log.Println(fmt.Sprintf("Got error logging user %s in: %s", formattedEmailAddress, err.Error()))
		if lErr != nil {
			writeJSONResponse(w, loginUserResponse{
				LoginError: lErr,
			})
			return
		}
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	redirectKeyForLocation := routes.GetLoginRedirectKeyOrDefault(req.RedirectKey)
	redirectURL, err := routes.GetLoginRedirectRouteForKeyAndUser(redirectKeyForLocation, *userID)
	if err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	if err := middleware.AssignAuthToken(w, *userID); err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	// This is a hack
	redirectPath := fmt.Sprintf("/%s", strings.TrimPrefix(*redirectURL, env.GetAbsoluteURLForEnvironment("")))
	writeJSONResponse(w, loginUserResponse{
		Location: ptr.String(redirectPath),
	})
}

type createUserRequest struct {
	CreateUserToken string `json:"create_user_token"`
	EmailAddress    string `json:"email_address"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type createUserResponse struct {
	CheckoutToken   *string          `json:"checkout_token"`
	CreateUserError *createUserError `json:"create_user_error"`
}

type createUserError string

const (
	createUserErrorAlreadyExists        createUserError = "already-exists"
	createUserErrorInvalidToken         createUserError = "invalid-token"
	createUserErrorPasswordRequirements createUserError = "pass-requirements"
	createUserErrorNoSubscription       createUserError = "no-subscription"
	createUserErrorPasswordsNoMatch     createUserError = "passwords-no-match"
)

func (c createUserError) Ptr() *createUserError {
	return &c
}

func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	var req createUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	userID, err := routetoken.ValidateTokenAndEmailAndGetUserID(req.CreateUserToken, routes.CreateUserKey, formattedEmailAddress)
	if err != nil {
		writeJSONResponse(w, createUserResponse{
			CreateUserError: createUserErrorInvalidToken.Ptr(),
		})
		return
	}
	if req.Password != req.ConfirmPassword {
		writeJSONResponse(w, createUserResponse{
			CreateUserError: createUserErrorPasswordsNoMatch.Ptr(),
		})
		return
	}
	if !useraccounts.ValidatePasswordMeetsRequirements(req.Password) {
		writeJSONResponse(w, createUserResponse{
			CreateUserError: createUserErrorPasswordRequirements.Ptr(),
		})
		return
	}
	var cErr *createUserError
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		alreadyHasAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
		switch {
		case err != nil:
			return err
		case alreadyHasAccount:
			cErr = createUserErrorAlreadyExists.Ptr()
			return fmt.Errorf("user already has account")
		}
		holdUntilTime := time.Now().Add(30 * time.Minute)
		if _, err := useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypeAccountCreated, holdUntilTime); err != nil {
			return err
		}
		if err := useraccounts.CreateUserPasswordForUser(tx, *userID, req.Password); err != nil {
			return err
		}
		_, err = bgstripe.CreateCustomerForUser(tx, *userID)
		return err
	}); err != nil {
		log.Println(fmt.Sprintf("Error signing up user %s: %s", formattedEmailAddress, err.Error()))
		if cErr != nil {
			writeJSONResponse(w, createUserResponse{
				CreateUserError: cErr,
			})
			return
		}
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	token, err := routes.MakePremiumSubscriptionCheckoutToken(*userID)
	if err != nil {
		log.Println("Error here")
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	if err := middleware.AssignAuthToken(w, *userID); err != nil {
		log.Println("Error here 2")
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	writeJSONResponse(w, createUserResponse{
		CheckoutToken: token,
	})
}

type getUserProfileRequest struct {
	SubscriptionManagementToken string `json:"subscription_management_token"`
}

type getUserProfileResponse struct {
	EmailAddress      *string                         `json:"email_address,omitempty"`
	SubscriptionLevel *useraccounts.SubscriptionLevel `json:"subscription_level,omitempty"`
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Got error getting user profile: %s", err.Error()))
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	var req getUserProfileRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println(fmt.Sprintf("Got error getting user profile: %s", err.Error()))
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	expectedUserID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		log.Println(fmt.Sprintf("Got error getting user profile: %s", err.Error()))
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	middleware.WithAuthorizationCheck(w, r, middleware.WithAuthorizationCheckInput{
		HandleFoundUser: func(userID users.UserID, subscriptionLevel *useraccounts.SubscriptionLevel, w http.ResponseWriter, r *http.Request) {
			if *expectedUserID != userID {
				middleware.RemoveAuthToken(w)
				writeJSONResponse(w, getUserProfileResponse{})
				return
			}
			var user *users.User
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				user, err = users.GetUser(tx, userID)
				return err
			}); err != nil {
				writeErrorJSONResponse(w, errorResponse{
					Message: "Request is not valid",
				})
				return
			}
			writeJSONResponse(w, getUserProfileResponse{
				EmailAddress:      ptr.String(user.EmailAddress),
				SubscriptionLevel: subscriptionLevel,
			})
		},
		HandleNoUserFound: func(w http.ResponseWriter, r *http.Request) {
			writeJSONResponse(w, getUserProfileResponse{})
		},
		HandleInvalidAuthenticationToken: func(w http.ResponseWriter, r *http.Request) {
			writeJSONResponse(w, getUserProfileResponse{})
		},
		HandleError: middleware.HandleAuthorizationError,
	})
}

type resetPasswordRequest struct {
	ResetPasswordToken string `json:"reset_password_token"`
	EmailAddress       string `json:"email_address"`
	Password           string `json:"password"`
	ConfirmPassword    string `json:"confirm_password"`
}

type resetPasswordResponse struct {
	ManagementToken    *string             `json:"management_token"`
	ResetPasswordError *resetPasswordError `json:"reset_password_error"`
}

type resetPasswordError string

const (
	resetPasswordErrorInvalidToken         resetPasswordError = "invalid-token"
	resetPasswordErrorTokenExpired         resetPasswordError = "token-expired"
	resetPasswordErrorPasswordRequirements resetPasswordError = "pass-requirements"
	resetPasswordErrorPasswordsNoMatch     resetPasswordError = "passwords-no-match"
	resetPasswordErrorNoAccount            resetPasswordError = "no-account"
)

func (r resetPasswordError) Ptr() *resetPasswordError {
	return &r
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	var req resetPasswordRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	var rErr *resetPasswordError
	var forgotPasswordID *useraccounts.ForgotPasswordAttemptID
	if err := encrypt.WithDecodedToken(req.ResetPasswordToken, func(tokenPair encrypt.TokenPair) error {
		if tokenPair.Key != routes.ForgotPasswordKey.Str() {
			rErr = resetPasswordErrorInvalidToken.Ptr()
			return fmt.Errorf("Invalid token")
		}
		forgotPasswordIDStr, ok := tokenPair.Value.(string)
		if !ok {
			return fmt.Errorf("Bad value for token pair")
		}
		forgotPasswordID = useraccounts.ForgotPasswordAttemptID(forgotPasswordIDStr).Ptr()
		return nil
	}); err != nil {
		log.Println(fmt.Sprintf("Error resetting password for user %s: %s", formattedEmailAddress, err.Error()))
		if rErr != nil {
			writeJSONResponse(w, resetPasswordResponse{
				ResetPasswordError: rErr,
			})
			return
		}
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	if req.Password != req.ConfirmPassword {
		writeJSONResponse(w, resetPasswordResponse{
			ResetPasswordError: resetPasswordErrorPasswordsNoMatch.Ptr(),
		})
		return
	}
	if !useraccounts.ValidatePasswordMeetsRequirements(req.Password) {
		writeJSONResponse(w, resetPasswordResponse{
			ResetPasswordError: resetPasswordErrorPasswordRequirements.Ptr(),
		})
		return
	}
	var userID *users.UserID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		forgotPasswordAttempt, isExpired, err := useraccounts.GetUnexpiredForgotPasswordAttemptByID(tx, *forgotPasswordID)
		switch {
		case err != nil:
			return err
		case isExpired:
			rErr = resetPasswordErrorTokenExpired.Ptr()
			return fmt.Errorf("Token has expired")
		case forgotPasswordAttempt == nil:
			return fmt.Errorf("Token has not expired, no error, but no forgot password attempt")
		}
		user, err := users.LookupUserForIDAndEmail(tx, forgotPasswordAttempt.UserID, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case user == nil:
			rErr = resetPasswordErrorInvalidToken.Ptr()
			return fmt.Errorf("The user does not correspond to the token")
		}
		userID = &user.ID
		alreadyHasAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, user.ID)
		switch {
		case err != nil:
			return err
		case !alreadyHasAccount:
			rErr = resetPasswordErrorNoAccount.Ptr()
			return fmt.Errorf("user does not have account")
		}
		if err := useraccounts.CreateUserPasswordForUser(tx, user.ID, req.Password); err != nil {
			return err
		}
		return useraccounts.SetForgotPasswordAttemptAsUsed(tx, forgotPasswordAttempt.ID)
	}); err != nil {
		log.Println(fmt.Sprintf("Error resetting password for user %s: %s", formattedEmailAddress, err.Error()))
		if rErr != nil {
			writeJSONResponse(w, resetPasswordResponse{
				ResetPasswordError: rErr,
			})
			return
		}
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	if err := middleware.AssignAuthToken(w, *userID); err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	token, err := routes.MakeSubscriptionManagementToken(*userID)
	if err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	writeJSONResponse(w, resetPasswordResponse{
		ManagementToken: token,
	})
}
