package emailtemplates

import "babblegraph/model/email"

const (
	genericEmailTemplateFilename        string = "generic_email_with_optional_action_link_template.html"
	genericNonUserEmailTemplateFilename string = "generic_non_user_template.html"
)

type genericTemplate struct {
	EmailTitle       string
	PreheaderText    string
	BeforeParagraphs []string
	Action           *GenericEmailAction
	AfterParagraphs  []string
}

type genericEmailWithOptionalActionTemplate struct {
	BaseEmailTemplate
	genericTemplate
}

type GenericEmailAction struct {
	Link       string
	ButtonText string
}

type MakeGenericUserEmailHTMLInput struct {
	EmailRecordID      email.ID
	UserAccessor       UserAccessor
	EmailTitle         string
	PreheaderText      string
	BeforeParagraphs   []string
	GenericEmailAction *GenericEmailAction
	AfterParagraphs    []string
}

func MakeGenericUserEmailHTML(input MakeGenericUserEmailHTMLInput) (*string, error) {
	baseEmailTemplate, err := createBaseEmailTemplate(input.EmailRecordID, input.UserAccessor)
	if err != nil {
		return nil, err
	}
	return openAndExecuteTemplate(genericEmailTemplateFilename, genericEmailWithOptionalActionTemplate{
		BaseEmailTemplate: *baseEmailTemplate,
		genericTemplate: genericTemplate{
			EmailTitle:       input.EmailTitle,
			PreheaderText:    input.PreheaderText,
			BeforeParagraphs: input.BeforeParagraphs,
			Action:           input.GenericEmailAction,
			AfterParagraphs:  input.AfterParagraphs,
		},
	})
}

type MakeGenericEmailHTMLInput struct {
	EmailTitle         string
	PreheaderText      string
	BeforeParagraphs   []string
	GenericEmailAction *GenericEmailAction
	AfterParagraphs    []string
}

func MakeGenericEmailHTML(input MakeGenericEmailHTMLInput) (*string, error) {
	return openAndExecuteTemplate(genericNonUserEmailTemplateFilename, genericTemplate{
		EmailTitle:       input.EmailTitle,
		PreheaderText:    input.PreheaderText,
		BeforeParagraphs: input.BeforeParagraphs,
		Action:           input.GenericEmailAction,
		AfterParagraphs:  input.AfterParagraphs,
	})
}
