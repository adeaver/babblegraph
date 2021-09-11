package emailtemplates

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
	"os"
)

func getPathForTemplateFile(filename string) (*string, error) {
	templatePath := env.GetEnvironmentVariableOrDefault("EMAIL_TEMPLATES_PATH", "/model/emailtemplates/templates/")
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ptr.String(fmt.Sprintf("%s%s%s", cwd, templatePath, filename)), nil
}

func createBaseEmailTemplate(emailRecordID email.ID, userAccessor UserAccessor) (*BaseEmailTemplate, error) {
	var subscriptionManagementLink *string
	if userAccessor.doesUserAlreadyHaveAccount() {
		subscriptionManagementLink = ptr.String(routes.GetLoginRoute())
	} else {
		var err error
		subscriptionManagementLink, err = routes.MakeSubscriptionManagementRouteForUserID(userAccessor.getUserID())
		if err != nil {
			return nil, err
		}
	}
	heroImageURL, err := routes.MakeLogoURLForEmailRecordID(emailRecordID)
	if err != nil {
		return nil, err
	}
	return &BaseEmailTemplate{
		SubscriptionManagementLink: *subscriptionManagementLink,
		HeroImageURL:               *heroImageURL,
		HomePageURL:                routes.MustGetHomePageURL(),
	}, nil
}
