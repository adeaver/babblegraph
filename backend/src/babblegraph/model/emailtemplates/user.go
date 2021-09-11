package emailtemplates

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"

	"github.com/jmoiron/sqlx"
)

type UserAccessor interface {
	getUserID() users.UserID
	doesUserAlreadyHaveAccount() bool
}

type DefaultUserAccessor struct {
	userID         users.UserID
	userHasAccount bool
}

func GetDefaultUserAccessor(tx *sqlx.Tx, userID users.UserID) (*DefaultUserAccessor, error) {
	userHasAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, userID)
	if err != nil {
		return nil, err
	}
	return &DefaultUserAccessor{
		userID:         userID,
		userHasAccount: userHasAccount,
	}, nil
}

func (d *DefaultUserAccessor) getUserID() users.UserID {
	return d.userID
}

func (d *DefaultUserAccessor) doesUserAlreadyHaveAccount() bool {
	return d.userHasAccount
}
