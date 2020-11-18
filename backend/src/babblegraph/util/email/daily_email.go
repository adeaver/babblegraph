package email

import (
	"babblegraph/util/ptr"
	"fmt"
	"strings"
	"text/template"
	"time"
)

type DailyEmailLink struct {
	ImageURL *string
	Title    *string
	URL      string
}

type dailyEmailTemplate struct {
	Links []DailyEmailLink
}

func (cl *Client) SendDailyEmailForLinks(recipient string, links []DailyEmailLink) error {
	emailBody, err := createEmailBody(links)
	if err != nil {
		return err
	}
	today := time.Now()
	return cl.sendEmail(sendEmailInput{
		Recipient: recipient,
		HTMLBody:  *emailBody,
		Subject:   fmt.Sprintf("Babblegraph Daily Links - %s %d, %d", today.Month().String(), today.Day(), today.Year()),
	})
}

func createEmailBody(links []DailyEmailLink) (*string, error) {
	t, err := template.New("Daily Email").ParseFiles("./templates/daily_email_template.html")
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, dailyEmailTemplate{
		Links: links,
	}); err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
