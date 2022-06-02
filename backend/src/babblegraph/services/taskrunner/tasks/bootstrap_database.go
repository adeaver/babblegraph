package tasks

import (
	"babblegraph/model/content"
	"babblegraph/services/taskrunner/bootstrap"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func BootstrapDatabase() error {
	switch env.MustEnvironmentName() {
	case env.EnvironmentProd,
		env.EnvironmentStage:
		return fmt.Errorf("Can only bootstrap in local environment")
	case env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail,
		env.EnvironmentTest:
		return database.WithTx(func(tx *sqlx.Tx) error {
			topicIDsByLabel := make(map[string]content.TopicID)
			for topic, displayName := range bootstrap.ContentTopics {
				topicID, err := content.AddTopic(tx, topic, true)
				if err != nil {
					return err
				}
				topicIDsByLabel[topic] = *topicID
				_, err = content.AddTopicDisplayName(tx, *topicID, wordsmith.LanguageCodeSpanish, displayName, true)
				if err != nil {
					return err
				}
			}
			for domain, info := range bootstrap.Sources {
				sourceID, err := content.InsertSource(tx, content.InsertSourceInput{
					Title:                 domain,
					LanguageCode:          info.LanguageCode,
					URL:                   domain,
					Type:                  content.SourceTypeNewsWebsite,
					IngestStrategy:        content.IngestStrategyWebsiteHTML1,
					Country:               info.Country,
					ShouldUseURLAsSeedURL: true,
					IsActive:              true,
				})
				if err != nil {
					return err
				}
				for _, u := range info.SeedURLs {
					parsed := urlparser.MustParseURL(u.URL)
					sourceSeedID, err := content.AddSourceSeed(tx, *sourceID, parsed, true)
					if err != nil {
						return err
					}
					topicID, ok := topicIDsByLabel[u.TopicLabel]
					if !ok {
						ctx.GetDefaultLogContext().Infof("No topic ID for label %s", topicID)
						continue
					}
					if err := content.UpsertSourceSeedMapping(tx, *sourceSeedID, topicID, true); err != nil {
						return err
					}
				}
			}
			return nil
		})
	default:
		return fmt.Errorf("Unrecognized environment")
	}
}
