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
	upsertExperimentQuery    = `INSERT INTO
        experiments (
            name,
            current_step,
            previous_step,
            is_active
        ) VALUES (
            $1, $2, $3, $4
        ) ON CONFLICT (name) DO UPDATE SET
            current_step=$2,
            previous_step=$3,
            is_active=$4,
            last_modified_at = timezone('utc', now())`

	getUserExperimentVariationQuery    = "SELECT * FROM experiments_user_variations WHERE experiment_id = $1 AND user = $2"
	insertUserExperimentVariationQuery = `
    INSERT INTO
        experiments_user_variations (
            experiment_id, user_id, in_variation, in_experiment, accessed_at_step
        ) VALUES (
            $1, $2, $3, $4, $5
        ) ON CONFLICT (experiment_id, user_id) DO UPDATE SET
            in_variation = $3,
            in_experiment = $4,
            accessed_at_step = $5,
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

func SetCurrentStepForExperiment(tx *sqlx.Tx, experimentName string, nextStep decimal.Number, defaultActive bool) error {
	experiment, err := lookupExperimentByName(tx, experimentName)
	switch {
	case err != nil:
		return err
	case experiment == nil:
		return upsertExperiment(tx, &dbExperiment{
			Name:         experimentName,
			CurrentStep:  nextStep.ToInt64Rounded(),
			PreviousStep: 0,
			IsActive:     defaultActive,
		})
	default:
		experiment.PreviousStep = experiment.CurrentStep
		experiment.CurrentStep = nextStep.ToInt64Rounded()
		return upsertExperiment(tx, experiment)
	}
}

func SetExperimentIsActive(tx *sqlx.Tx, experimentName string, isActive bool) error {
	experiment, err := lookupExperimentByName(tx, experimentName)
	switch {
	case err != nil:
		return err
	case experiment == nil:
		return upsertExperiment(tx, &dbExperiment{
			Name:         experimentName,
			CurrentStep:  0,
			PreviousStep: 0,
			IsActive:     isActive,
		})
	default:
		experiment.IsActive = isActive
		return upsertExperiment(tx, experiment)
	}
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
	_, err := tx.Exec(upsertExperimentQuery, experiment.Name, experiment.CurrentStep, experiment.PreviousStep, experiment.IsActive)
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
			var inExperiment, inVariation bool
			experimentRoll := rand.Int63n(100)
			inExperiment = experimentRoll < experiment.CurrentStep*2
			if inExperiment {
				roll := rand.Int63n(100)
				inVariation = roll < 50
			}
			if _, err := tx.Exec(insertUserExperimentVariationQuery, experiment.ID, userID, inVariation, inExperiment, experiment.CurrentStep); err != nil {
				return false, err
			}
			return inExperiment && inVariation, nil
		case len(matches) == 1:
			isUserInVariation := matches[0].InVariation && (experiment.IsActive || overrideInactiveExperiment)
			if !isUserInVariation && experiment.IsActive && matches[0].AccessedAtStep != experiment.CurrentStep {
				var isUserInExperiment bool
				experimentProbability := decimal.FromInt64(experiment.CurrentStep * 2).Divide(decimal.FromInt64(100 - experiment.PreviousStep*2))
				experimentRoll := rand.Int63n(100)
				isUserInExperiment = experimentRoll < experimentProbability.ToInt64Rounded()
				if isUserInExperiment {
					roll := rand.Int63n(100)
					isUserInVariation = roll < 50
				}
				if _, err := tx.Exec(insertUserExperimentVariationQuery, experiment.ID, userID, isUserInVariation, isUserInExperiment, experiment.CurrentStep); err != nil {
					return false, err
				}
			}
			return isUserInVariation, nil
		}
	}
	return false, nil
}
