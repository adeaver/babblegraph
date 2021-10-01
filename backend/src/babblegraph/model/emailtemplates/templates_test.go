package emailtemplates

import (
	"babblegraph/model/email"
	"babblegraph/model/newsletter"
	"babblegraph/model/users"
	"babblegraph/util/ptr"
	"strings"
	"testing"
)

func TestCreateNewsletterTemplate(t *testing.T) {
	userAccessor := &testUserAccessor{
		userHasAccount: false,
		userID:         users.UserID("12345"),
	}
	testNewsletter := newsletter.NewsletterBody{
		SetTopicsLink:     ptr.String("babblegraph.com/topics"),
		ReinforcementLink: "babblegraph.com/reinforce",
	}
	html, err := MakeNewsletterHTML(MakeNewsletterHTMLInput{
		EmailRecordID: email.NewEmailRecordID(),
		UserAccessor:  userAccessor,
		Body:          testNewsletter,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(*html) == 0 {
		t.Fatalf("Got empty body")
	}
	if !strings.Contains(*html, "babblegraph.com/topics") {
		t.Errorf("Expected topics link")
	}
	if !strings.Contains(*html, "babblegraph.com/reinforce") {
		t.Errorf("Expected reinforcement link")
	}
}
