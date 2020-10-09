package sendutil

import (
	"babblegraph/model/documents"
	"babblegraph/util/env"
	"fmt"
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
		body := makeEmailBody(docs)
		smtpAuth := smtp.PlainAuth("", creds.from, creds.pass, creds.smtpServerHost)
		if err := smtp.SendMail(creds.GetSMTPServerWithPort(), smtpAuth, creds.from, []string{emailAddress}, []byte(body)); err != nil {
			return err
		}
	}
	return nil
}

func makeEmailBody(docs []documents.Document) string {
	var urls []string
	for _, doc := range docs {
		urls = append(urls, doc.URL)
	}
	return strings.Join(urls, "\n")
}
