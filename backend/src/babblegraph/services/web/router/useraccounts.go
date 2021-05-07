package router

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/middleware"
	"babblegraph/services/web/util/routetoken"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func registerUserAccountsRoutes() {
	a.prefixes["useraccounts"] = true
	a.r.HandleFunc("/api/useraccounts/login_user_1", middleware.WithoutBodyLogger(loginUser))
	a.routeNames["/api/useraccounts/login_user_1"] = true
	a.r.HandleFunc("/api/useraccounts/create_user_1", middleware.WithoutBodyLogger(createUser))
	a.routeNames["/api/useraccounts/create_user_1"] = true
}

type loginUserRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

type loginUserResponse struct {
	ManagementToken *string     `json:"management_token"`
	LoginError      *loginError `json:"login_error"`
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
	token, err := routes.MakeSubscriptionManagementToken(*userID)
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
	writeJSONResponse(w, loginUserResponse{
		ManagementToken: token,
	})
}

type createUserRequest struct {
	CreateUserToken string `json:"create_user_token"`
	EmailAddress    string `json:"email_address"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type createUserResponse struct {
	ManagementToken *string          `json:"management_token"`
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
	// TODO: validate password
	var cErr *createUserError
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
		switch {
		case err != nil:
			return err
		case subscriptionLevel == nil:
			cErr = createUserErrorNoSubscription.Ptr()
			return fmt.Errorf("no subscription found for user")
		}
		// TODO: validate that user doesn't already have account
		return useraccounts.CreateUserPasswordForUser(tx, *userID, req.Password)
	}); err != nil {
		log.Println(fmt.Sprintf("Error signing up user %s: %s", formattedEmailAddress, err.Error()))
		if cErr != nil {
			writeJSONResponse(w, createUserResponse{
				CreateUserError: createUserErrorInvalidToken.Ptr(),
			})
			return
		}
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
	if err := middleware.AssignAuthToken(w, *userID); err != nil {
		writeErrorJSONResponse(w, errorResponse{
			Message: "Request is not valid",
		})
		return
	}
	writeJSONResponse(w, createUserResponse{
		ManagementToken: token,
	})
}
