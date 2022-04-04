package experiment

import (
	"babblegraph/model/users"
	"time"
)

type experimentID string

type dbExperiment struct {
	CreatedAt      time.Time    `db:"created_at"`
	LastModifiedAt time.Time    `db:"last_modified_at"`
	ID             experimentID `db:"_id"`
	Name           string       `db:"name"`
	CurrentStep    int64        `db:"current_step"`
	PreviousStep   int64        `db:"previous_step"`
	IsActive       bool         `db:"is_active"`
}

type experimentUserVariationID string

type dbExperimentUserVariation struct {
	CreatedAt      time.Time                 `db:"created_at"`
	LastModifiedAt time.Time                 `db:"last_modified_at"`
	ID             experimentUserVariationID `db:"_id"`
	UserID         users.UserID              `db:"user_id"`
	ExperimentID   experimentID              `db:"experiment_id"`
	InExperiment   bool                      `db:"in_experiment"`
	InVariation    bool                      `db:"in_variation"`
	AccessedAtStep int64                     `db:"accessed_at_step"`
}
