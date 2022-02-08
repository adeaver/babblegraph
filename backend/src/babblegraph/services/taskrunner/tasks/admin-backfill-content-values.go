package tasks

import (
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

func BackfillAdminContentValues() error {
	return database.WithTx(func(tx *sqlx.Tx) error {

	})
}
