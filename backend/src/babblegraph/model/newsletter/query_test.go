package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/podcasts"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"babblegraph/wordsmith"
	"strings"
	"testing"
)

func TestUserHasAccount(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		PodcastAccessor:   &testPodcastAccessor{},
		ContentAccessor:   &testContentAccessor{},
	})
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
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		ContentAccessor:   &testContentAccessor{},
	})
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
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []content.TopicID{
			content.TopicID("topicid-art"),
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		PodcastAccessor:   &testPodcastAccessor{},
		ContentAccessor:   &testContentAccessor{},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	body := testNewsletter.Body
	if err := testutils.CompareNullableString(body.SetTopicsLink, nil); err != nil {
		t.Errorf("Error on set topics link: %s", err.Error())
	}
	if body.ReinforcementLink != routes.MakeLoginLinkWithReinforcementRedirect() {
		t.Errorf("Error on reinforcement link. Expected %s, but got %s", routes.MakeLoginLinkWithReinforcementRedirect(), body.ReinforcementLink)
	}
}

func TestUserScheduleDay(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []content.TopicID{
			content.TopicID("topicid-art"),
		},
		userSubscriptionLevel: useraccounts.SubscriptionLevelPremium.Ptr(),
		userNewsletterSchedule: usernewsletterschedule.TestNewsletterSchedule{
			SendRequested: false,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		PodcastAccessor:   &testPodcastAccessor{},
		ContentAccessor:   &testContentAccessor{},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if testNewsletter != nil {
		t.Errorf("Expected null newsletter, but it was not")
	}
}

func TestUserScheduleDayNoSubscription(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []content.TopicID{
			content.TopicID("topicid-art"),
		},
		userNewsletterSchedule: usernewsletterschedule.TestNewsletterSchedule{
			SendRequested: false,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		PodcastAccessor:   &testPodcastAccessor{},
		ContentAccessor:   &testContentAccessor{},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if testNewsletter == nil {
		t.Errorf("Expected non-null newsletter, but it was not")
	}
}

func TestWholeNewsletterHasPodcasts(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	emailRecordID := email.NewEmailRecordID()
	userAccessor := &testUserAccessor{
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userTopics: []content.TopicID{
			content.TopicID("test-art"),
			content.TopicID("test-astronomy"),
			content.TopicID("test-architecture"),
			content.TopicID("test-automotive"),
		},
		allowableSourceIDs: []content.SourceID{
			content.SourceID("test-source"),
		},
		doesUserHaveAccount: true,
	}
	documentTopics := []content.TopicID{
		content.TopicID("test-art"),
		content.TopicID("test-astronomy"),
		content.TopicID("test-architecture"),
		content.TopicID("test-automotive"),
		content.TopicID("test-culture"),
	}
	contentAccessor := &testContentAccessor{}
	var docs []documents.DocumentWithScore
	var podcasts []podcasts.Episode
	for idx, topic := range documentTopics {
		doc, _, err := getDefaultDocumentWithLink(c, idx, emailRecordID, contentAccessor, userAccessor, getDefaultDocumentInput{
			Topics: []content.TopicID{topic},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		docs = append(docs, *doc)
		podcasts = append(podcasts, getDefaultPodcast(topic))
	}
	podcastAccessor := &testPodcastAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
		validSourceIDs: []content.SourceID{
			content.SourceID("test-source"),
		},
		podcastEpisodes: podcasts,
	}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor: &testDocsAccessor{
			documents: docs,
		},
		PodcastAccessor: podcastAccessor,
		ContentAccessor: &testContentAccessor{},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if testNewsletter == nil {
		t.Errorf("Expected non-null newsletter, but it was not")
	}
	if len(testNewsletter.Body.Categories) == 0 {
		t.Errorf("Expected newsletter body with categories, but it was not")
	}
	for _, c := range testNewsletter.Body.Categories {
		if len(c.PodcastLinks) != 0 {
			t.Errorf("Error on category with name %s: expected 0 podcast, but got %d", *c.Name, len(c.PodcastLinks))
		}
	}
}
