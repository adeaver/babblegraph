package emailtemplates

import "babblegraph/model/users"

type testUserAccessor struct {
	userID         users.UserID
	userHasAccount bool
}

func (d *testUserAccessor) getUserID() users.UserID {
	return d.userID
}

func (d *testUserAccessor) doesUserAlreadyHaveAccount() bool {
	return d.userHasAccount
}
