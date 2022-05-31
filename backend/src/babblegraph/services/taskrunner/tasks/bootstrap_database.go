package tasks

import (
	"babblegraph/model/content"
	"babblegraph/services/taskrunner/bootstrap"
	"babblegraph/util/database"
	"babblegraph/util/env"
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
			for topic, displayName := range bootstrap.ContentTopics {
				topicID, err := content.AddTopic(tx, topic, true)
				if err != nil {
					return err
				}
				_, err = content.AddTopicDisplayName(tx, *topicID, wordsmith.LanguageCodeSpanish, displayName, true)
				if err != nil {
					return err
				}
			}
			return nil
		})
	default:
		return fmt.Errorf("Unrecognized environment")
	}
}
