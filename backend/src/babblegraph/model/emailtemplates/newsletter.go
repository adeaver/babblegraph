package emailtemplates

import (
	"babblegraph/model/email"
	"babblegraph/model/newsletter"
	"babblegraph/util/ptr"
	"html/template"
	"strings"
)

const (
	newsletterTemplateFilename = "newsletter_template.html"
)

type newsletterTemplate struct {
	BaseEmailTemplate
	Body newsletter.NewsletterBody
}

type MakeNewsletterHTMLInput struct {
	EmailRecordID email.ID
	UserAccessor  UserAccessor
	Body          newsletter.NewsletterBody
}

func MakeNewsletterHTML(input MakeNewsletterHTMLInput) (*string, error) {
	baseEmailTemplate, err := createBaseEmailTemplate(input.EmailRecordID, input.UserAccessor)
	if err != nil {
		return nil, err
	}
	templateFile, err := getPathForTemplateFile(newsletterTemplateFilename)
	if err != nil {
		return nil, err
	}
	t, err := template.New(newsletterTemplateFilename).ParseFiles(*templateFile)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, newsletterTemplate{
		BaseEmailTemplate: *baseEmailTemplate,
		Body:              input.Body,
	}); err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
