package sendutil

import (
	"babblegraph/model/documents"
	"babblegraph/util/env"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var sentFromEmail string

func InitializeEmailClient() error {
	sentFromEmail = env.MustEnvironmentVariable("EMAIL_ADDRESS")
	return nil
}

func SendEmailsToUser(emailAddressesToDocuments map[string][]documents.Document) error {
	awsAccessKey := env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY")
	awsSecretAccessKey := env.MustEnvironmentVariable("AWS_SES_SECRET_KEY")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretAccessKey, ""),
	})
	if err != nil {
		return err
	}
	svc := ses.New(sess)

	// Assemble the email.
	for emailAddress, docs := range emailAddressesToDocuments {
		if len(docs) == 0 {
			log.Println(fmt.Sprintf("No docs for user %s, skipping", emailAddress))
		}
		body, err := makeEmailBody(emailAddress, docs)
		if err != nil {
			return nil
		}
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(emailAddress),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Text: &ses.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(*body),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String("Babblegraph Daily Links"),
				},
			},
			Source: aws.String(sentFromEmail),
		}

		// Attempt to send the email.
		if _, err := svc.SendEmail(input); err != nil {
			return err
		}
		log.Println(fmt.Sprintf("Sent an email to %s", emailAddress))
	}
	return nil
}

var emailTemplate = `{{range $val := .URLs}}{{$val}}
{{end}}`

type emailBodyInfo struct {
	URLs []string
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
		URLs: urls,
	}); err != nil {
		return nil, err
	}
	body := b.String()
	return &body, nil
}
