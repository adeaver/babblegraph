package email

import "babblegraph/model/users"

type Recipient struct {
	EmailAddress string
	UserID       users.UserID
}

// All email templates should use this
type BaseEmailTemplate struct {
	SubscriptionManagementLink string
	UnsubscribeLink            string
	HeroImageURL               string
	HomePageURL                string
}
