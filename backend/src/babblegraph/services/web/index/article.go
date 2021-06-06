package index

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlinks"
	"babblegraph/model/users"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/util/random"
	"babblegraph/util/urlparser"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func getArticleLinkBodyFromMap(m map[string]interface{}) (*routes.ArticleLinkBodyDEPRECATED, error) {
	var articleBody routes.ArticleLinkBodyDEPRECATED
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Got error %s marshalling map %+v", err.Error(), m)
	}
	if err := json.Unmarshal(bytes, &articleBody); err != nil {
		return nil, fmt.Errorf("Got error %s unmarshalling string: %s", err.Error(), string(bytes))
	}
	return &articleBody, nil
}

func HandleArticleLink(w http.ResponseWriter, r *http.Request) {
	router.LogRequestWithoutBody(r)
	routeVars := mux.Vars(r)
	token, ok := routeVars["token"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("handle-article-link", fmt.Sprintf("article-link-%s", random.MustMakeRandomString(12)))
	})
	var emailRecordID *email.ID
	var userID *users.UserID
	var url *urlparser.ParsedURL
	if err := encrypt.WithDecodedToken(token, func(tokenPair encrypt.TokenPair) error {
		switch {
		case tokenPair.Key == routes.ArticleLinkKeyDEPRECATED.Str():
			m, ok := tokenPair.Value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Article body did not marshal correctly, got type %v", reflect.TypeOf(tokenPair.Value))
			}
			articleBody, err := getArticleLinkBodyFromMap(m)
			if err != nil {
				return err
			}
			emailRecordID = &articleBody.EmailRecordID
			userID = &articleBody.UserID
			u := urlparser.MustParseURL(articleBody.URL)
			url = &u
			return nil
		case tokenPair.Key == routes.ArticleLinkKeyForUserDocumentID.Str():
			userDocumentIDStr, ok := tokenPair.Value.(string)
			if !ok {
				return fmt.Errorf("Article body did not marshal correctly, got type %v", reflect.TypeOf(tokenPair.Value))
			}
			userDocumentID := userdocuments.UserDocumentID(userDocumentIDStr)
			return database.WithTx(func(tx *sqlx.Tx) error {
				userDocument, err := userdocuments.GetUserDocumentID(tx, userDocumentID)
				if err != nil {
					return err
				}
				if userDocument.DocumentURL == nil {
					return fmt.Errorf("User Document has no document URL")
				}
				emailRecordID = userDocument.EmailID
				userID = &userDocument.UserID
				u := urlparser.MustParseURL(*userDocument.DocumentURL)
				url = &u
				return nil
			})
		default:
			return fmt.Errorf("Incorrect key type: %s", tokenPair.Key)
		}
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	if emailRecordID == nil || userID == nil || url == nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(fmt.Errorf("Token does not have email record, user id, or url: %s", token))
		return
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return userlinks.RegisterUserLinkClick(tx, *userID, *url, *emailRecordID)
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	http.Redirect(w, r, url.URL, http.StatusFound)
}

func getPaywallReportBodyFromMap(m map[string]interface{}) (*routes.PaywallReportBodyDEPRECATED, error) {
	var paywallReportBody routes.PaywallReportBodyDEPRECATED
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Got error %s marshalling map %+v", err.Error(), m)
	}
	if err := json.Unmarshal(bytes, &paywallReportBody); err != nil {
		return nil, fmt.Errorf("Got error %s unmarshalling string: %s", err.Error(), string(bytes))
	}
	return &paywallReportBody, nil
}

func HandlePaywallReport(w http.ResponseWriter, r *http.Request) {
	router.LogRequestWithoutBody(r)
	routeVars := mux.Vars(r)
	token, ok := routeVars["token"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("handle-paywall-report-link", fmt.Sprintf("paywall-report-link-%s", random.MustMakeRandomString(12)))
	})
	var emailRecordID *email.ID
	var userID *users.UserID
	var url *urlparser.ParsedURL
	if err := encrypt.WithDecodedToken(token, func(tokenPair encrypt.TokenPair) error {
		switch {
		case tokenPair.Key == routes.PaywallReportKeyDEPRECATED.Str():
			m, ok := tokenPair.Value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Paywall report body did not marshal correctly, got type %v", reflect.TypeOf(tokenPair.Value))
			}
			paywallReportBody, err := getPaywallReportBodyFromMap(m)
			if err != nil {
				return err
			}
			emailRecordID = &paywallReportBody.EmailRecordID
			userID = &paywallReportBody.UserID
			u := urlparser.MustParseURL(paywallReportBody.URL)
			url = &u
			return nil
		case tokenPair.Key == routes.PaywallReportKeyForUserDocumentID.Str():
			userDocumentIDStr, ok := tokenPair.Value.(string)
			if !ok {
				return fmt.Errorf("Paywall report body did not marshal correctly, got type %v", reflect.TypeOf(tokenPair.Value))
			}
			userDocumentID := userdocuments.UserDocumentID(userDocumentIDStr)
			return database.WithTx(func(tx *sqlx.Tx) error {
				userDocument, err := userdocuments.GetUserDocumentID(tx, userDocumentID)
				if err != nil {
					return err
				}
				if userDocument.DocumentURL == nil {
					return fmt.Errorf("User Document has no document URL")
				}
				emailRecordID = userDocument.EmailID
				userID = &userDocument.UserID
				u := urlparser.MustParseURL(*userDocument.DocumentURL)
				url = &u
				return nil
			})
		default:
			return fmt.Errorf("Incorrect key type: %s", tokenPair.Key)
		}
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	if emailRecordID == nil || userID == nil || url == nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(fmt.Errorf("Token does not have email record, user id, or url: %s", token))
		return
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return userlinks.ReportPaywall(tx, *userID, *url, *emailRecordID)
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	subscriptionManagementLink, err := routes.MakeSubscriptionManagementToken(*userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/paywall-thank-you/%s", *subscriptionManagementLink), http.StatusPermanentRedirect)
}
