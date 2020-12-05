package email

import "babblegraph/model/users"

type Recipient struct {
	EmailAddress string
	UserID       users.UserID
}
