package experiment

import (
	"babblegraph/model/users"
	"babblegraph/util/math/decimal"
	"fmt"
	"math/rand"

	"github.com/jmoiron/sqlx"
)

const (
	getExperimentByNameQuery = "SELECT * FROM experiments WHERE name=$1"
	updateExperimentQuery    = "UPDATE experiments SET current_step = $1, previous_step = $2, is_active = $3, last_modified_at = timezone('utc', now()) WHERE _id = $4"

	getUserExperimentVariationQuery    = "SELECT * FROM experiments_user_variations WHERE experiment_id = $1 AND user = $2"
	insertUserExperimentVariationQuery = `
    INSERT INTO
        experiments_user_variations (
            experiment_id, user_id, in_variation, accessed_at_step
        ) VALUES (
            $1, $2, $3, $4
        ) ON CONFLICT (experiment_id, user_id) DO UPDATE SET
            in_variation = $3,
            accessed_at_step = $4,
            last_modified_at = timezone('utc', now())
    `
)

func GetCurrentStepForExperiment(tx *sqlx.Tx, experimentName string) (*decimal.Number, error) {
	experiment, err := lookupExperimentByName(tx, experimentName)
	switch {
	case err != nil:
		return nil, err
	case experiment == nil:
		return nil, fmt.Errorf("No experiment found for name %s", experimentName)
	default:
		return decimal.FromInt64(experiment.CurrentStep).Ptr(), nil
	}
}

func SetCurrentStepForExperiment(tx *sqlx.Tx, experimentName string, nextStep decimal.Number) error {
	return nil
}

func SetExperimentIsActive(tx *sqlx.Tx, experimentName string, isActive bool) error {
	return nil
}

func lookupExperimentByName(tx *sqlx.Tx, experimentName string) (*dbExperiment, error) {
	var matches []dbExperiment
	err := tx.Select(&matches, getExperimentByNameQuery, experimentName)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most 1 experiment with the name %s but got %d", experimentName, len(matches))
	default:
		m := matches[0]
		return &m, nil
	}
}

func upsertExperiment(tx *sqlx.Tx, experiment *dbExperiment) error {
	_, err := tx.Exec(updateExperimentQuery, experiment.IsActive, experiment.CurrentStep, experiment.PreviousStep, experiment.ID)
	return err
}

func IsUserInVariation(tx *sqlx.Tx, experimentName string, userID users.UserID, overrideInactiveExperiment bool) (bool, error) {
	experiment, err := lookupExperimentByName(tx, experimentName)
	switch {
	case err != nil:
		return false, err
	case experiment == nil:
		return false, fmt.Errorf("No experiment found for name %s", experimentName)
	default:
		var matches []dbExperimentUserVariation
		err := tx.Select(&matches, getUserExperimentVariationQuery, experiment.ID, userID)
		switch {
		case err != nil:
			return false, err
		case len(matches) > 1:
			return false, fmt.Errorf("Expected at most one user variation entry for ID %s, but got %d", userID, len(matches))
		case len(matches) == 0 && experiment.IsActive:
			roll := rand.Int63n(100)
			inVariation := roll < experiment.CurrentStep
			if _, err := tx.Exec(insertUserExperimentVariationQuery, experiment.ID, userID, inVariation, experiment.CurrentStep); err != nil {
				return false, err
			}
			return inVariation, nil
		case len(matches) == 1:
			isUserInVariation := matches[0].InVariation && (experiment.IsActive || overrideInactiveExperiment)
			if !isUserInVariation && experiment.IsActive && matches[0].AccessedAtStep != experiment.CurrentStep {
				probability := decimal.FromInt64(experiment.CurrentStep).Divide(decimal.FromInt64(100 - experiment.PreviousStep))
				roll := rand.Int63n(100)
				isUserInVariation = roll < probability.ToInt64Rounded()
				if _, err := tx.Exec(insertUserExperimentVariationQuery, experiment.ID, userID, isUserInVariation, experiment.CurrentStep); err != nil {
					return false, err
				}
			}
			return isUserInVariation, nil
		}
	}
	return false, nil
}
