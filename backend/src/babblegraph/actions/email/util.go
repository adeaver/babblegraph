package email

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func getPathForTemplateFile(filename string) (*string, error) {
	templatePath := env.GetEnvironmentVariableOrDefault("TEMPLATES_PATH", "/actions/email/templates/")
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ptr.String(fmt.Sprintf("%s%s%s", cwd, templatePath, filename)), nil
}

func createBaseTemplate(tx *sqlx.Tx, emailRecordID email.ID, recipient email.Recipient) (*email.BaseEmailTemplate, error) {
	var subscriptionManagementLink *string
	alreadyHasAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, recipient.UserID)
	if err != nil {
		return nil, err
	}
	if alreadyHasAccount {
		subscriptionManagementLink = ptr.String(env.GetAbsoluteURLForEnvironment("login"))
	} else {
		subscriptionManagementLink, err = routes.MakeSubscriptionManagementRouteForUserID(recipient.UserID)
		if err != nil {
			return nil, err
		}
	}
	heroImageURL, err := routes.MakeLogoURLForEmailRecordID(emailRecordID)
	if err != nil {
		return nil, err
	}
	return &email.BaseEmailTemplate{
		SubscriptionManagementLink: *subscriptionManagementLink,
		HeroImageURL:               *heroImageURL,
		HomePageURL:                routes.MustGetHomePageURL(),
	}, nil
}
