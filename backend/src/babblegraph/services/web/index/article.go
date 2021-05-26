package index

import (
	"babblegraph/model/routes"
	"babblegraph/model/userlinks"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/util/urlparser"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func HandleArticleLink(w http.ResponseWriter, r *http.Request) {
	router.LogRequestWithoutBody(r)
	routeVars := mux.Vars(r)
	token, ok := routeVars["token"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var articleBody routes.ArticleLinkBody
	if err := encrypt.WithDecodedToken(token, func(tokenPair encrypt.TokenPair) error {
		if tokenPair.Key != routes.ArticleLinkKey.Str() {
			return fmt.Errorf("Wrong key type")
		}
		var ok bool
		articleBody, ok = tokenPair.Value.(routes.ArticleLinkBody)
		if !ok {
			return fmt.Errorf("Article Body unmarshalled incorrectly")
		}
		return nil
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		u := urlparser.MustParseURL(articleBody.URL)
		return userlinks.RegisterUserLinkClick(tx, articleBody.UserID, u, articleBody.EmailRecordID)
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, articleBody.URL, http.StatusPermanentRedirect)
}
