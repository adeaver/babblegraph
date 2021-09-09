package newsletter

import (
	"babblegraph/model/email"
	"babblegraph/model/users"

	"github.com/jmoiron/sqlx"
)

type emailAccessor interface {
	InsertEmailRecord(id email.ID, userID users.UserID) error
}

type DefaultEmailAccessor struct {
	tx *sqlx.Tx
}

func CreateDefaultEmailAccessor(tx *sqlx.Tx) *DefaultEmailAccessor {
	return &DefaultEmailAccessor{tx: tx}
}

func (d *DefaultEmailAccessor) InsertEmailRecord(id email.ID, userID users.UserID) error {
	return email.InsertEmailRecord(d.tx, id, userID, email.EmailTypeDaily)
}
