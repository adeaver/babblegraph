package documents

import (
	"babblegraph/util/urlparser"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func makeDocumentIndexForURL(parsedURL urlparser.ParsedURL) DocumentID {
	md5Hash := md5.Sum([]byte(parsedURL.URLIdentifier))
	return DocumentID(fmt.Sprintf("web_doc-%s", hex.EncodeToString(md5Hash[:])))
}
