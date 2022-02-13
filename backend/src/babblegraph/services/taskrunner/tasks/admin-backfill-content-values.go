package tasks

import (
	"babblegraph/model/content"
	"babblegraph/model/links2"
	"babblegraph/model/urltopicmapping"
	"babblegraph/model/usercontenttopics"
	"babblegraph/model/userlinks"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

func BackfillAdminContentValues() error {
	c := ctx.GetDefaultLogContext()
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return usercontenttopics.BackfillUserContentTopicMappings(c, tx)
	}); err != nil {
		return err
	}
	c.Infof("Starting user link clicks backfill")
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return userlinks.BackfillUserLinkClicks(c, tx)
	}); err != nil {
		return err
	}
	c.Infof("Finished user link clicks backfill")
	c.Infof("Starting user newsletter schedule backfill")
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return usernewsletterschedule.BackfillNewsletterScheduleTopics(c, tx)
	}); err != nil {
		return err
	}
	c.Infof("Finished user newsletter schedule backfill")
	return database.WithTx(func(tx *sqlx.Tx) error {
		var count int
		return links2.GetLinksCursor(tx, func(link links2.Link) (bool, error) {
			return false, database.WithTx(func(tx *sqlx.Tx) error {
				parsedURL := urlparser.ParseURL(link.URL)
				if parsedURL == nil {
					c.Infof("URL %s did not parse correctly", link.URL)
					return nil
				}
				sourceID, err := content.GetSourceIDForParsedURL(tx, *parsedURL)
				switch {
				case err != nil:
					c.Infof("Error getting source ID for URL %s", link.URL)
					return err
				case sourceID == nil:
					c.Infof("URL %s did not produce a source ID", link.URL)
					return nil
				}
				count++
				if count%1000 == 0 {
					c.Infof("Successfully processed %d links", count)
				}
				// Update link
				if err := links2.InsertLinks(tx, []urlparser.ParsedURL{
					*parsedURL,
				}); err != nil {
					return err
				}
				topics, mappings, err := urltopicmapping.GetTopicsAndMappingIDsForURL(tx, link.URL)
				switch {
				case err != nil:
					return err
				case len(topics) == len(mappings):
					// No need to update
				default:
					var topicMappings []urltopicmapping.TopicMappingUnion
					for _, t := range topics {
						topicID, err := content.GetTopicIDByContentTopic(tx, t)
						if err != nil {
							return err
						}
						topicMappingID, err := content.LookupTopicMappingIDForSourceAndTopic(c, tx, *sourceID, *topicID)
						switch {
						case err != nil:
							return err
						case topicMappingID == nil:
							c.Infof("URL %s and topic %s produced null topic mapping id, skipping", link.URL, t)
						default:
							topicMappings = append(topicMappings, urltopicmapping.TopicMappingUnion{
								Topic:          t,
								TopicMappingID: *topicMappingID,
							})
						}
					}
					return urltopicmapping.ApplyContentTopicsToURL(tx, *parsedURL, topicMappings)
				}
				return nil
			})
		})
	})
}
