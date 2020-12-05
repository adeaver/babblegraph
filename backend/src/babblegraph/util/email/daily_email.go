package email

import (
	"babblegraph/util/ptr"
	"babblegraph/util/routes"
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
	UnsubscribeLink string
	Links           []DailyEmailLink
}

func (cl *Client) SendDailyEmailForLinks(recipient Recipient, links []DailyEmailLink) error {
	for _, l := range links {
		if l.Description != nil {
			log.Println(fmt.Sprintf("Email util description found %s for URL %s", *l.Description, l.URL))
		} else {
			log.Println(fmt.Sprintf("Email util no description found for URL %s", l.URL))
		}
	}
	emailBody, err := createEmailBody(dailyEmailTemplate{
		UnsubscribeLink: routes.MakeUnsubscribeRouteForUserID(recipient.UserID),
		Links:           links,
	})
	if err != nil {
		return err
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
