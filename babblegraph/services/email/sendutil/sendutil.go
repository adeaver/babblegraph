package sendutil

import (
	"babblegraph/model/documents"
	"babblegraph/util/env"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"
)

var creds *emailCredentials

type emailCredentials struct {
	from           string
	password       string
	smtpServerHost string
	smtpServerPort string
}

func (c emailCredentials) GetSMTPServerWithPort() string {
	return fmt.Sprintf("%s:%s", c.smtpServerHost, c.smtpServerPort)
}

func InitializeEmailClient() error {
	creds = &emailCredentials{
		from:           env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		password:       env.MustEnvironmentVariable("EMAIL_PASSWORD"),
		smtpServerHost: env.GetEnvironmentVariableOrDefault("SMTP_SERVER_HOST", "smtp.gmail.com"),
		smtpServerPort: env.GetEnvironmentVariableOrDefault("SMTP_SERVER_PORT", "587"),
	}
	return nil
}

func SendEmailsToUser(emailAddressesToDocuments map[string][]documents.Document) error {
	if creds == nil {
		panic("email client not configured")
	}
	for emailAddress, docs := range emailAddressesToDocuments {
		if len(docs) == 0 {
			log.Println(fmt.Sprintf("No docs for user %s, skipping", emailAddress))
		}
		body, err := makeEmailBody(emailAddress, docs)
		if err != nil {
			return nil
		}
		smtpAuth := smtp.PlainAuth("", creds.from, creds.password, creds.smtpServerHost)
		if err := smtp.SendMail(creds.GetSMTPServerWithPort(), smtpAuth, creds.from, []string{emailAddress}, []byte(*body)); err != nil {
			return err
		}
		log.Println(fmt.Sprintf("Sent an email to %s", emailAddress))
	}
	return nil
}

var emailTemplate = `To: {{.Recipient}}
Subject: {{.Subject}}

{{range $val := .URLs}}
{{$val}}{{end}}`

type emailBodyInfo struct {
	Recipient string
	Subject   string
	URLs      []string
}

func makeEmailBody(recipient string, docs []documents.Document) (*string, error) {
	bodyTemplate := template.Must(template.New("email").Parse(emailTemplate))
	var urls []string
	for _, doc := range docs {
		urls = append(urls, doc.URL)
	}
	log.Println(fmt.Sprintf("Sending the following URLS %+v to %s", urls, recipient))
	var b strings.Builder
	if err := bodyTemplate.Execute(&b, emailBodyInfo{
		Recipient: recipient,
		Subject:   "Babblegraph Daily Links",
		URLs:      urls,
	}); err != nil {
		return nil, err
	}
	body := b.String()
	return &body, nil
}
