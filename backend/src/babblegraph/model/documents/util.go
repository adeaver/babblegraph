package documents

import (
	"babblegraph/util/urlparser"
	"crypto/md5"
	"encoding/hex"
)

func makeDocumentIndexForURL(parsedURL urlparser.ParsedURL) DocumentID {
	md5Hash := md5.Sum([]byte(parsedURL.URLIdentifier))
	return DocumentID(hex.EncodeToString(md5Hash[:]))
}
