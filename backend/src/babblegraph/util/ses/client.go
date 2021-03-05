package ses

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Client struct {
	awsAccessKey       string
	awsSecretAccessKey string
	awsRegion          string
	fromAddress        string
}

type NewClientInput struct {
	AWSAccessKey       string
	AWSSecretAccessKey string
	AWSRegion          string
	FromAddress        string
}

func NewClient(input NewClientInput) *Client {
	return &Client{
		awsAccessKey:       input.AWSAccessKey,
		awsSecretAccessKey: input.AWSSecretAccessKey,
		awsRegion:          input.AWSRegion,
		fromAddress:        input.FromAddress,
	}
}

type SendEmailInput struct {
	Recipient string
	HTMLBody  string
	Subject   string
}

func (cl *Client) SendEmail(input SendEmailInput) (*string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cl.awsRegion),
		Credentials: credentials.NewStaticCredentials(cl.awsAccessKey, cl.awsSecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	svc := ses.New(sess)
	sesInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(input.Recipient),
			},
			CcAddresses: []*string{},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(input.HTMLBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(input.Subject),
			},
		},
		Source: aws.String(fmt.Sprintf("\"Babblegraph\" <%s>", cl.fromAddress)),
	}
	output, err := svc.SendEmail(sesInput)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("Sent email with id %s to %s", *output.MessageId, input.Recipient))
	return output.MessageId, nil
}
