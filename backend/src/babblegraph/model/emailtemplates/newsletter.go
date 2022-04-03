package emailtemplates

import (
	"babblegraph/model/email"
	"babblegraph/model/newsletter"
	"babblegraph/util/ptr"
	"html/template"
	"strings"
)

const (
	newsletterTemplateFilename         = "newsletter_template.html"
	newsletterTemplateVersion2Filename = "version2_newsletter_template.html"
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

type newsletterVersion2Template struct {
	BaseEmailTemplate
	Body newsletter.NewsletterVersion2Body
}

type MakeNewsletterVersion2HTMLInput struct {
	EmailRecordID email.ID
	UserAccessor  UserAccessor
	Body          newsletter.NewsletterVersion2Body
}

func MakeNewsletterVersion2HTML(input MakeNewsletterVersion2HTMLInput) (*string, error) {
	baseEmailTemplate, err := createBaseEmailTemplate(input.EmailRecordID, input.UserAccessor)
	if err != nil {
		return nil, err
	}
	return openAndExecuteTemplate(newsletterTemplateVersion2Filename, newsletterVersion2Template{
		BaseEmailTemplate: *baseEmailTemplate,
		Body:              input.Body,
	})
}
