package email

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/encrypt"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func HandleDailyEmailOpenToken(tx *sqlx.Tx, token string) error {
	return encrypt.WithDecodedToken(token, func(t encrypt.TokenPair) error {
		if t.Key != routes.EmailOpenedKey.Str() {
			return fmt.Errorf("Token has wrong key: %s", t.Key)
		}
		emailRecordID, ok := t.Value.(string)
		if !ok {
			return fmt.Errorf("Token has wrong value type")
		}
		return email.SetEmailFirstOpened(tx, email.ID(emailRecordID))
	})
}
