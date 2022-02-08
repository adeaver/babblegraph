package tasks

import (
	"babblegraph/model/documents"
	"babblegraph/model/links2"
	"babblegraph/model/urltopicmapping"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

func ReindexDocuments() error {
	c := ctx.GetDefaultLogContext()
	return database.WithTx(func(tx *sqlx.Tx) error {
		var count int64
		return links2.GetLinksCursor(tx, func(link links2.Link) (bool, error) {
			switch {
			case link.LastFetchVersion != nil && *link.LastFetchVersion == links2.FetchVersion8:
				c.Infof("URL %s already has current version, skipping", link.URL)
				return false, nil
			case link.SourceID == nil:
				c.Infof("URL %s has a null source ID, skipping", link.URL)
				return false, nil
			}
			count++
			if count%1000 == 0 {
				c.Infof("Processed %d links", count)
			}
			parsedURL := urlparser.ParseURL(link.URL)
			if parsedURL == nil {
				c.Infof("URL %s did not parse correctly", link.URL)
				return false, nil
			}
			return false, database.WithTx(func(tx *sqlx.Tx) error {
				_, mappingIDs, err := urltopicmapping.GetTopicsAndMappingIDsForURL(tx, link.URL)
				switch {
				case err != nil:
					return err
				case len(mappingIDs) == 0:
					c.Infof("No mapping IDs for URL %s", link.URL)
				}
				return documents.UpdateDocumentForURL(*parsedURL, documents.UpdateDocumentInput{
					Version:         documents.Version8,
					TopicMappingIDs: mappingIDs,
					SourceID:        *link.SourceID,
				})
			})
		})
	})
}
