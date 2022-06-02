package article

import (
	"babblegraph/model/routes"
	"babblegraph/model/userdocuments"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"babblegraph/util/urlparser"
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "article",
	Routes: []router.Route{
		{
			Path: "get_article_metadata_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				getArticleMetadata,
			),
		},
	},
}

type getArticleMetadataRequest struct {
	ArticleToken string `json:"article_token"`
}

type getArticleMetadataResponse struct {
	ReaderToken string `json:"reader_token"`
	ArticleID   string `json:"article_id"`
}

func getArticleMetadata(r *router.Request) (interface{}, error) {
	var req getArticleMetadataRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var articleID string
	var userID *users.UserID
	if err := encrypt.WithDecodedToken(req.ArticleToken, func(tokenPair encrypt.TokenPair) error {
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
				userID = &userDocument.UserID
				if userDocument.DocumentURL == nil {
					return fmt.Errorf("User Document has no document URL")
				}
				u := urlparser.MustParseURL(*userDocument.DocumentURL)
				articleID = base64.URLEncoding.EncodeToString([]byte(u.URLIdentifier))
				return nil
			})
		default:
			return fmt.Errorf("Incorrect key type: %s", tokenPair.Key)
		}
	}); err != nil {
		return nil, err
	}
	readerToken, err := routes.EncryptUserIDWithKey(*userID, routes.ArticleReaderKey)
	if err != nil {
		return nil, err
	}
	return getArticleMetadataResponse{
		ArticleID:   articleID,
		ReaderToken: *readerToken,
	}, nil
}
