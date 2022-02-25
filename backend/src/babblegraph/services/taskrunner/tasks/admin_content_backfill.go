package tasks

import (
	"babblegraph/model/content"
	"babblegraph/model/domains"
	"babblegraph/model/links2"
	"babblegraph/model/urltopicmapping"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

func BackfillContent() error {
	c := ctx.GetDefaultLogContext()
	return database.WithTx(func(tx *sqlx.Tx) error {
		var successCount, errorCount, skippedCount int
		err := links2.GetLinksCursor(tx, func(link links2.Link) (bool, error) {
			u := urlparser.ParseURL(link.URL)
			skippedCount++
			switch {
			case u == nil:
				c.Debugf("URL %s did not parse correctly", link.URL)
				return false, nil
			case link.SourceID != nil:
				c.Debugf("Source is already present")
				return false, nil
			case u == nil:
				c.Debugf("URL did not parse correctly")
				return false, nil
			case !domains.IsURLAllowed(*u):
				c.Debugf("Domain is not allowed, skipping")
				return false, nil
			}
			skippedCount--
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				sourceID, err := link.GetSourceID(tx)
				if err != nil {
					return err
				}
				if err := links2.UpdateLinkSource(tx, *u, *sourceID); err != nil {
					return err
				}
				topics, _, err := urltopicmapping.GetTopicsAndMappingIDsForURL(tx, link.URL)
				if err != nil {
					return err
				}
				if len(topics) == 0 {
					return nil
				}
				var topicMappings []urltopicmapping.TopicMappingUnion
				for _, t := range topics {
					topicID, err := content.GetTopicIDByContentTopic(tx, t)
					if err != nil {
						return err
					}
					topicMappingID, err := content.LookupTopicMappingIDForSourceAndTopic(c, tx, *sourceID, *topicID)
					if err != nil {
						return err
					}
					topicMappings = append(topicMappings, urltopicmapping.TopicMappingUnion{
						Topic:          t,
						TopicMappingID: *topicMappingID,
					})
				}
				return urltopicmapping.ApplyContentTopicsToURL(tx, *u, topicMappings)
			}); err != nil {
				c.Errorf("Error on url %s: %s", link.URL, err.Error())
				errorCount++
				return false, nil
			}
			successCount++
			if (successCount+errorCount+skippedCount)%1000 == 0 {
				c.Infof("Successfully processed %d links", (successCount + errorCount + skippedCount))
			}
			return false, nil
		})
		c.Infof("Success %d, Error %d, Skipped %d", successCount, errorCount, skippedCount)
		return err
	})
}
