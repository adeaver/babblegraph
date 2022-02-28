package tasks

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/links2"
	"babblegraph/model/urltopicmapping"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ReindexDocuments() error {
	c := ctx.GetDefaultLogContext()
	return database.WithTx(func(tx *sqlx.Tx) error {
		var count int64
		return links2.GetLinksCursor(tx, func(link links2.Link) (bool, error) {
			parsedURL := urlparser.ParseURL(link.URL)
			if parsedURL == nil {
				c.Infof("URL %s did not parse correctly", link.URL)
				return false, nil
			}
			return false, database.WithTx(func(tx *sqlx.Tx) error {
				sourceID, err := link.GetSourceID(tx)
				switch {
				case err != nil:
					return err
				case sourceID == nil:
					return fmt.Errorf("Link with URL %s does not have source ID, skipping", link.URL)
				}
				var topicIDs []content.TopicID
				_, mappingIDs, err := urltopicmapping.GetTopicsAndMappingIDsForURL(tx, link.URL)
				switch {
				case err != nil:
					return err
				case len(mappingIDs) == 0:
					c.Infof("No mapping IDs for URL %s", link.URL)
				default:
					var sourceSeedTopicMappingIDs []content.SourceSeedTopicMappingID
					for _, mappingID := range mappingIDs {
						sourceSeedTopicMapping, sourceTopicMapping, err := mappingID.GetOriginID()
						switch {
						case err != nil:
							c.Errorf("Error getting origin ID %s", err.Error())
						case sourceTopicMapping != nil:
							c.Errorf("Unsupported ID type: source topic mapping for url: %s", parsedURL.URLIdentifier)
						case sourceSeedTopicMapping != nil:
							sourceSeedTopicMappingIDs = append(sourceSeedTopicMappingIDs, *sourceSeedTopicMapping)
						}
					}
					if len(sourceSeedTopicMappingIDs) > 0 {
						topicIDs, err = content.LookupTopicsForSourceSeedMappingIDs(tx, sourceSeedTopicMappingIDs)
						if err != nil {
							return err
						}
					}
				}
				count++
				if count%1000 == 0 {
					c.Infof("Processed %d links", count)
				}
				return documents.UpdateDocumentForURL(*parsedURL, documents.UpdateDocumentInput{
					Version:         documents.Version8,
					TopicMappingIDs: mappingIDs,
					TopicIDs:        topicIDs,
					SourceID:        *sourceID,
				})
			})
		})
	})
}
