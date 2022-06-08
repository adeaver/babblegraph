package index

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/userdocuments"
	"babblegraph/model/userlinks"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/middleware"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/util/urlparser"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/html"
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
		routermiddleware.RemoveUnsafeCookies(w, r)
		url, err := getURLForUserDocument(r)
		switch {
		case err != nil:
			http.Error(w, http.StatusText(500), 500)
			return
		case url == nil:
			http.Error(w, http.StatusText(404), 404)
			return
		}
		u := urlparser.MustParseURL(*url)
		resp, err := http.Get(*url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		htmlDoc, err := html.Parse(strings.NewReader(string(data)))
		if err != nil {
			log.Fatalf(err.Error())
		}
		var f func(node *html.Node)
		f = func(node *html.Node) {
			switch node.Data {
			case "script",
				"link",
				"style":
				var attrs []html.Attribute
				for _, attr := range node.Attr {
					if (attr.Key == "href" || attr.Key == "src") && strings.HasPrefix(attr.Val, "/") {
						absoluteURL, err := urlparser.EnsureProtocol(fmt.Sprintf("%s/%s", u.Domain, attr.Val))
						if err != nil {
							continue
						}
						attr.Val = *absoluteURL
					}
					attrs = append(attrs, attr)
				}
				node.Attr = attrs
			}
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(htmlDoc)
		var b bytes.Buffer
		err = html.Render(&b, htmlDoc)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Write([]byte(b.String()))
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
