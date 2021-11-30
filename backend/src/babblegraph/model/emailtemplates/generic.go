package emailtemplates

import "babblegraph/model/email"

const (
	genericEmailTemplateFilename string = "generic_email_with_optional_action_link_template.html"
)

type genericEmailWithOptionalActionTemplate struct {
	BaseEmailTemplate
	EmailTitle       string
	PreheaderText    string
	BeforeParagraphs []string
	Action           *GenericEmailAction
	AfterParagraphs  []string
}

type GenericEmailAction struct {
	Link       string
	ButtonText string
}

type MakeGenericEmailHTMLInput struct {
	EmailRecordID      email.ID
	UserAccessor       UserAccessor
	EmailTitle         string
	PreheaderText      string
	BeforeParagraphs   []string
	GenericEmailAction *GenericEmailAction
	AfterParagraphs    []string
}

func MakeGenericEmailHTML(input MakeGenericEmailHTMLInput) (*string, error) {
	baseEmailTemplate, err := createBaseEmailTemplate(input.EmailRecordID, input.UserAccessor)
	if err != nil {
		return nil, err
	}
	return openAndExecuteTemplate(genericEmailTemplateFilename, genericEmailWithOptionalActionTemplate{
		BaseEmailTemplate: *baseEmailTemplate,
		EmailTitle:        input.EmailTitle,
		PreheaderText:     input.PreheaderText,
		BeforeParagraphs:  input.BeforeParagraphs,
		Action:            input.GenericEmailAction,
		AfterParagraphs:   input.AfterParagraphs,
	})
}
