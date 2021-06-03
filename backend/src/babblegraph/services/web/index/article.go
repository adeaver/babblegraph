package index

import (
	"babblegraph/model/routes"
	"babblegraph/model/userlinks"
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

func getArticleLinkBodyFromMap(m map[string]interface{}) (*routes.ArticleLinkBody, error) {
	var articleBody routes.ArticleLinkBody
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
	var articleBody *routes.ArticleLinkBody
	if err := encrypt.WithDecodedToken(token, func(tokenPair encrypt.TokenPair) error {
		if tokenPair.Key != routes.ArticleLinkKey.Str() {
			return fmt.Errorf("Wrong key type")
		}
		m, ok := tokenPair.Value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Article body did not marshal correctly, got type %v", reflect.TypeOf(tokenPair.Value))
		}
		var err error
		articleBody, err = getArticleLinkBodyFromMap(m)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		u := urlparser.MustParseURL(articleBody.URL)
		return userlinks.RegisterUserLinkClick(tx, articleBody.UserID, u, articleBody.EmailRecordID)
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	http.Redirect(w, r, articleBody.URL, http.StatusFound)
}

func getPaywallReportBodyFromMap(m map[string]interface{}) (*routes.PaywallReportBody, error) {
	var paywallReportBody routes.PaywallReportBody
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &paywallReportBody); err != nil {
		return nil, err
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
	var paywallReportBody *routes.PaywallReportBody
	if err := encrypt.WithDecodedToken(token, func(tokenPair encrypt.TokenPair) error {
		if tokenPair.Key != routes.PaywallReportKey.Str() {
			return fmt.Errorf("Wrong key type")
		}
		var ok bool
		m, ok := tokenPair.Value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Article Body unmarshalled incorrectly")
		}
		var err error
		paywallReportBody, err = getPaywallReportBodyFromMap(m)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		u := urlparser.MustParseURL(paywallReportBody.URL)
		return userlinks.ReportPaywall(tx, paywallReportBody.UserID, u, paywallReportBody.EmailRecordID)
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	subscriptionManagementLink, err := routes.MakeSubscriptionManagementToken(paywallReportBody.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		localHub.CaptureException(err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/paywall-thank-you/%s", *subscriptionManagementLink), http.StatusPermanentRedirect)
}
