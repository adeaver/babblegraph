package scheduler

import (
	"babblegraph/model/useraccounts"
	"babblegraph/util/database"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func expireUserAccounts() error {
	var numExpiredAccounts *int64
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		numExpiredAccounts, err = useraccounts.SetExpiredSubscriptionsAsInactive(tx)
		return err
	}); err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Expired %d subscriptions", *numExpiredAccounts))
	return nil
}
