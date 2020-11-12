package sendutil

import (
	"babblegraph/model/documents"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

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
			continue
		}
		body, err := makeEmailBody(emailAddress, docs)
		if err != nil {
			return nil
		}
		today := time.Now()
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(emailAddress),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(*body),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(fmt.Sprintf("Babblegraph Daily Links - %s %d, %d", today.Month().String(), today.Day(), today.Year())),
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

var emailTemplate = `
<style>
.page-link {
    width: 100%;
    padding: 2%;
}

.page-image {
    display: block;
    margin: auto;
    max-height: 200px;
    min-height: 100px;
    width: auto;
}

.title {
    text-align: center;
}
</style>

{{range $page := .Links}}
<a class="page-link" href="{{$page.URL}}">
    {{if $page.ImageURL}}
        <img class="page-image" src="{{$page.ImageURL}}" />
        <p class="title">{{$page.Title}}</p>
    {{else if $page.Title}}
        <p class="title">[{{$page.Domain}}] - {{$page.Title}}</p>
    {{else}}
        <p class="url">{{$page.URL}}</p>
    {{end}}
</a>
<br />
{{end}}`

type emailBodyInfo struct {
	Links []linkInfo
}

type linkInfo struct {
	URL      string
	Domain   string
	Title    *string
	ImageURL *string
}

func makeEmailBody(recipient string, docs []documents.Document) (*string, error) {
	bodyTemplate := template.Must(template.New("email").Parse(emailTemplate))
	var links []linkInfo
	for _, doc := range docs {
		var title, imageURL *string
		if doc.Metadata != nil {
			if len(doc.Metadata.Title) > 0 {
				title = ptr.String(doc.Metadata.Title)
			}
			if len(doc.Metadata.Image) > 0 && urlparser.IsValidURL(doc.Metadata.Image) {
				imageURL = ptr.String(doc.Metadata.Image)
			}
		}
		log.Println("Sending URL %s to recipient %s", doc.URL, recipient)
		parsedURL := urlparser.ParseURL(doc.URL)
		if parsedURL == nil {
			continue
		}
		links = append(links, linkInfo{
			URL:      doc.URL,
			Domain:   parsedURL.Domain,
			ImageURL: imageURL,
			Title:    title,
		})
	}
	var b strings.Builder
	if err := bodyTemplate.Execute(&b, emailBodyInfo{
		Links: links,
	}); err != nil {
		return nil, err
	}
	body := b.String()
	return &body, nil
}
