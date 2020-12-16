package email

import (
	"babblegraph/model/routes"
	"babblegraph/util/ptr"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"
)

const dailyEmailTemplateFilename = "daily_email_template.html"

type DailyEmailLink struct {
	ImageURL    *string
	Title       *string
	Description *string
	URL         string
}

type dailyEmailTemplate struct {
	BaseEmailTemplate
	Links []DailyEmailLink
}

func (cl *Client) SendDailyEmailForLinks(recipient Recipient, links []DailyEmailLink) (*string, error) {
	for _, l := range links {
		if l.Description != nil {
			log.Println(fmt.Sprintf("Email util description found %s for URL %s", *l.Description, l.URL))
		} else {
			log.Println(fmt.Sprintf("Email util no description found for URL %s", l.URL))
		}
	}
	unsubscribeLink, err := routes.MakeUnsubscribeRouteForUserID(recipient.UserID)
	if err != nil {
		return nil, err
	}
	subscriptionManagementLink, err := routes.MakeSubscriptionManagementRouteForUserID(recipient.UserID)
	if err != nil {
		return nil, err
	}
	emailBody, err := createEmailBody(dailyEmailTemplate{
		BaseEmailTemplate: BaseEmailTemplate{
			SubscriptionManagementLink: *subscriptionManagementLink,
			UnsubscribeLink:            *unsubscribeLink,
		},
		Links: links,
	})
	if err != nil {
		return nil, err
	}
	today := time.Now()
	return cl.sendEmail(sendEmailInput{
		Recipient: recipient.EmailAddress,
		HTMLBody:  *emailBody,
		Subject:   fmt.Sprintf("Babblegraph Daily Links - %s %d, %d", today.Month().String(), today.Day(), today.Year()),
	})
}

func createEmailBody(emailData dailyEmailTemplate) (*string, error) {
	templateFile, err := getPathForTemplateFile(dailyEmailTemplateFilename)
	if err != nil {
		return nil, err
	}
	t, err := template.New(dailyEmailTemplateFilename).ParseFiles(*templateFile)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, emailData); err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
