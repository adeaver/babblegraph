package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"babblegraph/wordsmith"
	"strings"
	"testing"
)

func TestUserHasAccount(t *testing.T) {
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := &testEmailAccessor{}
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
	if err != nil {
		t.Fatalf(err.Error())
	}
	body := testNewsletter.Body
	if err := testutils.CompareNullableString(body.SetTopicsLink, ptr.String(routes.MakeLoginLinkWithContentTopicsRedirect())); err != nil {
		t.Errorf("Error on set topics link: %s", err.Error())
	}
	if body.ReinforcementLink != routes.MakeLoginLinkWithReinforcementRedirect() {
		t.Errorf("Error on reinforcement link. Expected %s, but got %s", routes.MakeLoginLinkWithReinforcementRedirect(), body.ReinforcementLink)
	}
}

func TestUserDoesNotHaveAccount(t *testing.T) {
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := &testEmailAccessor{}
	userAccessor := &testUserAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
	if err != nil {
		t.Fatalf(err.Error())
	}
	body := testNewsletter.Body
	if body.SetTopicsLink == nil || !strings.Contains(*body.SetTopicsLink, "/manage/") {
		t.Errorf("Error on set topics link. Expected it to contain /manage/ but got, but got %v", body.SetTopicsLink)
	}
	if !strings.Contains(body.ReinforcementLink, "/manage/") {
		t.Errorf("Error on set topics link. Expected it to contain /manage/ but got, but got %s", body.ReinforcementLink)
	}
}

func TestNoSetTopicsLink(t *testing.T) {
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := &testEmailAccessor{}
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
	if err != nil {
		t.Fatalf(err.Error())
	}
	body := testNewsletter.Body
	if err := testutils.CompareNullableString(body.SetTopicsLink, nil); err == nil {
		t.Errorf("Error on set topics link: %s", err.Error())
	}
	if body.ReinforcementLink != routes.MakeLoginLinkWithReinforcementRedirect() {
		t.Errorf("Error on reinforcement link. Expected %s, but got %s", routes.MakeLoginLinkWithReinforcementRedirect(), body.ReinforcementLink)
	}
}

func TestUserScheduleDay(t *testing.T) {
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := &testEmailAccessor{}
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
		userSubscriptionLevel: useraccounts.SubscriptionLevelPremium.Ptr(),
		userScheduleForDay: &usernewsletterschedule.UserNewsletterScheduleDayMetadata{
			IsActive: false,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if testNewsletter != nil {
		t.Errorf("Expected null newsletter, but it was not")
	}
}

func TestUserScheduleDayNoSubscription(t *testing.T) {
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := &testEmailAccessor{}
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
		userScheduleForDay: &usernewsletterschedule.UserNewsletterScheduleDayMetadata{
			IsActive: false,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if testNewsletter == nil {
		t.Errorf("Expected non-null newsletter, but it was not")
	}
}
