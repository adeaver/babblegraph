package index

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlinks"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/middleware"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/util/urlparser"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func getURLForUserDocument(r *http.Request) (*string, error) {
	routeVars := mux.Vars(r)
	token, ok := routeVars["token"]
	if !ok {
		return nil, nil
	}
	var url *string
	err := database.WithTx(func(tx *sqlx.Tx) error {
		userDocumentID := userdocuments.UserDocumentID(token)
		userDocument, err := userdocuments.GetUserDocumentID(tx, userDocumentID)
		if err != nil {
			return err
		}
		url = userDocument.DocumentURL
		return nil
	})
	switch {
	case err != nil:
		return nil, err
	case url == nil:
		return nil, nil
	}
	return url, nil
}

func handleArticleHTMLPassthrough() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := getURLForUserDocument(r)
		switch {
		case err != nil:
			http.Error(w, http.StatusText(500), 500)
			return
		case url == nil:
			http.Error(w, http.StatusText(404), 404)
			return
		}
		resp, err := http.Get(*url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		w.Write(data)
	}
}

func handleArticleOutLink() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := getURLForUserDocument(r)
		switch {
		case err != nil:
			http.Error(w, http.StatusText(500), 500)
			return
		case url == nil:
			http.Error(w, http.StatusText(404), 404)
			return
		}
		http.Redirect(w, r, *url, http.StatusFound)
	}
}
