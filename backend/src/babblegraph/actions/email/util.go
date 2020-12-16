package email

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
	"os"
)

func getPathForTemplateFile(filename string) (*string, error) {
	templatePath := env.GetEnvironmentVariableOrDefault("TEMPLATES_PATH", "/actions/email/templates/")
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ptr.String(fmt.Sprintf("%s%s%s", cwd, templatePath, filename)), nil
}

func createBaseTemplate(recipient email.Recipient) (*email.BaseEmailTemplate, error) {
	unsubscribeLink, err := routes.MakeUnsubscribeRouteForUserID(recipient.UserID)
	if err != nil {
		return nil, err
	}
	subscriptionManagementLink, err := routes.MakeSubscriptionManagementRouteForUserID(recipient.UserID)
	if err != nil {
		return nil, err
	}
	return &email.BaseEmailTemplate{
		SubscriptionManagementLink: *subscriptionManagementLink,
		UnsubscribeLink:            *unsubscribeLink,
	}, nil
}
