package index

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlinks"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

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

func handleArticleRoute(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.LogRequestWithoutBody(r)
		routeVars := mux.Vars(r)
		token, ok := routeVars["token"]
		if !ok {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		errs := make(chan error)
		async.WithContext(errs, "article-route", func(c async.Context) {
			var emailRecordID *email.ID
			var userID *users.UserID
			var url *urlparser.ParsedURL
			if err := encrypt.WithDecodedToken(token, func(tokenPair encrypt.TokenPair) error {
				switch {
				case tokenPair.Key == routes.ArticleLinkKeyDEPRECATED.Str():
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
				c.Infof("Unable to parse token: %s", err.Error())
				return
			}
			if emailRecordID == nil || userID == nil || url == nil {
				c.Warnf("Token does not have email record, user id, or url: %s", token)
				return
			}
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return userlinks.RegisterUserLinkClick(tx, *userID, *url, *emailRecordID)
			}); err != nil {
				c.Warnf("Failed to capture link click for user: %s", *userID)
			}
		}).Start()
		serveIndexTemplate(fmt.Sprintf("%s/index.html", staticFileDirName), w, r)
	}
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

func handlePaywallReport(r *router.Request) (interface{}, error) {
	token, err := r.GetRouteVar("token")
	if err != nil {
		return nil, err
	}
	var emailRecordID *email.ID
	var userID *users.UserID
	var url *urlparser.ParsedURL
	if err := encrypt.WithDecodedToken(*token, func(tokenPair encrypt.TokenPair) error {
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
		return nil, fmt.Errorf("Unable to parse token: %s", err.Error())
	}
	if emailRecordID == nil || userID == nil || url == nil {
		return nil, fmt.Errorf("Token does not have email record, user id, or url: %s", *token)
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return userlinks.ReportPaywall(tx, *userID, *url, *emailRecordID)
	}); err != nil {
		r.Warnf("Failed to capture paywall report for user: %s", *userID)
	}
	subscriptionManagementLink, err := routes.MakeSubscriptionManagementToken(*userID)
	if err != nil {
		return nil, err
	}
	return ptr.String(fmt.Sprintf("/paywall-thank-you/%s", *subscriptionManagementLink)), nil
}
