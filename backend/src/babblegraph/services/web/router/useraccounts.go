package router

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func registerUserAccountsRoutes() {
	a.prefixes["useraccounts"] = true
	a.r.HandleFunc("/api/useraccounts/login_user_1", withoutBodyLogger(loginUser))
	a.routeNames["/api/useraccounts/login_user_1"] = true
	a.r.HandleFunc("/api/useraccounts/create_user_1", withoutBodyLogger(createUser))
	a.routeNames["/api/useraccounts/create_user_1"] = true
}

type loginUserRequest struct {
	EmailAddress   string `json:"email_address"`
	Password       string `json:"password"`
	RedirectURLKey string `json:"redirect_url_key"`
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
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case user == nil:
			return fmt.Errorf("no user found for email address")
		}
		return useraccounts.VerifyPasswordForUser(tx, user.ID, req.Password)
	}); err != nil {
		// TODO: write a successful request with error message
	}
	// TODO: set JWT and redirect based on key
}

type createUserRequest struct {
	EmailAddress   string `json:"email_address"`
	Password       string `json:"password"`
	RedirectURLKey string `json:"redirect_url_key"`
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
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case user == nil:
			return fmt.Errorf("no user found for email address")
		}
		subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, user.ID)
		switch {
		case err != nil:
			return err
		case subscriptionLevel == nil:
			return fmt.Errorf("no subscription found for user")
		}
		return useraccounts.CreateUserPasswordForUser(tx, user.ID, req.Password)
	}); err != nil {
		// TODO: write a successful request with error message
	}
	// TODO: set JWT and redirect based on key
}
