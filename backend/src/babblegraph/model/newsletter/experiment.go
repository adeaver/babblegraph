package newsletter

import (
	"babblegraph/model/experiment"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/math/decimal"

	"github.com/jmoiron/sqlx"
)

const (
	version2ExperimentName       = "use_version_2_newsletter"
	version2ExperimentPercentage = 20
)

func InitializeNewsletterExperiments() error {
	return database.WithTx(func(tx *sqlx.Tx) error {
		return experiment.SetCurrentStepForExperiment(tx, version2ExperimentName, decimal.FromInt64(version2ExperimentPercentage), true)
	})
}

func IsUserInVersion2Experiment(tx *sqlx.Tx, userID users.UserID) (bool, error) {
	return experiment.IsUserInVariation(tx, version2ExperimentName, userID, false)
}
