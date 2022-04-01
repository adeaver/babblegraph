package newsletter

import (
	"babblegraph/model/advertising"
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/podcasts"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"babblegraph/wordsmith"
	"strings"
	"testing"
	"time"
)

func TestUserHasAccount(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor:     wordsmithAccessor,
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          docsAccessor,
		PodcastAccessor:       &testPodcastAccessor{},
		ContentAccessor:       &testContentAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{},
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
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor:     wordsmithAccessor,
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          docsAccessor,
		ContentAccessor:       &testContentAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{},
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
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor:     wordsmithAccessor,
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          docsAccessor,
		PodcastAccessor:       &testPodcastAccessor{},
		ContentAccessor:       &testContentAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{},
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
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested: false,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor:     wordsmithAccessor,
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          docsAccessor,
		PodcastAccessor:       &testPodcastAccessor{},
		ContentAccessor:       &testContentAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{},
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
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested: false,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor:     wordsmithAccessor,
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          docsAccessor,
		PodcastAccessor:       &testPodcastAccessor{},
		ContentAccessor:       &testContentAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if testNewsletter != nil {
		t.Errorf("Expected null newsletter, but it was not")
	}
}

func TestWholeNewsletterHasPodcasts(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	emailRecordID := email.NewEmailRecordID()
	userAccessor := &testUserAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
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
		doesUserHaveAccount:   true,
		userSubscriptionLevel: useraccounts.SubscriptionLevelPremium.Ptr(),
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
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
		PodcastAccessor:       podcastAccessor,
		ContentAccessor:       &testContentAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{},
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

func TestUserShouldShowAdvertisement(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userCreatedDate: time.Now().Add(-180 * 24 * time.Hour),
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		ContentAccessor:   &testContentAccessor{},
		PodcastAccessor:   &testPodcastAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{
			isUserEligibleForAdvertisement: true,
			advertisements: []advertising.Advertisement{
				{
					ID:           advertising.AdvertisementID("test-advertisement"),
					LanguageCode: wordsmith.LanguageCodeSpanish,
					CampaignID:   advertising.CampaignID("test-campaign"),
					Title:        "An Advertisement",
					ImageURL:     "www.babblegraph.com/test-image",
					Description:  "This is a description of an advertisement",
					IsActive:     true,
				},
			},
		},
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
	if body.Advertisement == nil {
		t.Errorf("Expected an advertisement, but there was none")
	}
}

func TestNoAdvertisementUserAccountAge(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userCreatedDate: time.Now().Add(-7 * 24 * time.Hour),
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		ContentAccessor:   &testContentAccessor{},
		PodcastAccessor:   &testPodcastAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{
			isUserEligibleForAdvertisement: true,
			advertisements: []advertising.Advertisement{
				{
					ID:           advertising.AdvertisementID("test-advertisement"),
					LanguageCode: wordsmith.LanguageCodeSpanish,
					CampaignID:   advertising.CampaignID("test-campaign"),
					Title:        "An Advertisement",
					ImageURL:     "www.babblegraph.com/test-image",
					Description:  "This is a description of an advertisement",
					IsActive:     true,
				},
			},
		},
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
	if body.Advertisement != nil {
		t.Errorf("Expected no advertisement, but there was one")
	}
}

func TestUserAdvertisementIneligible(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
		userCreatedDate: time.Now().Add(-180 * 24 * time.Hour),
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		ContentAccessor:   &testContentAccessor{},
		PodcastAccessor:   &testPodcastAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{
			isUserEligibleForAdvertisement: false,
			advertisements: []advertising.Advertisement{
				{
					ID:           advertising.AdvertisementID("test-advertisement"),
					LanguageCode: wordsmith.LanguageCodeSpanish,
					CampaignID:   advertising.CampaignID("test-campaign"),
					Title:        "An Advertisement",
					ImageURL:     "www.babblegraph.com/test-image",
					Description:  "This is a description of an advertisement",
					IsActive:     true,
				},
			},
		},
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
	if body.Advertisement != nil {
		t.Errorf("Expected no advertisement, but there was one")
	}
}

func TestUserSubscriptionHasNoAdvertisement(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		userSubscriptionLevel: useraccounts.SubscriptionLevelPremium.Ptr(),
		languageCode:          wordsmith.LanguageCodeSpanish,
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
		userCreatedDate: time.Now().Add(-180 * 24 * time.Hour),
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		ContentAccessor:   &testContentAccessor{},
		PodcastAccessor:   &testPodcastAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{
			isUserEligibleForAdvertisement: true,
			advertisements: []advertising.Advertisement{
				{
					ID:           advertising.AdvertisementID("test-advertisement"),
					LanguageCode: wordsmith.LanguageCodeSpanish,
					CampaignID:   advertising.CampaignID("test-campaign"),
					Title:        "An Advertisement",
					ImageURL:     "www.babblegraph.com/test-image",
					Description:  "This is a description of an advertisement",
					IsActive:     true,
				},
			},
		},
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
	if body.Advertisement != nil {
		t.Errorf("Expected no advertisement, but there was one")
	}
}

func TestNoAdvertisementIsOkay(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
		userCreatedDate: time.Now().Add(-180 * 24 * time.Hour),
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
		ContentAccessor:   &testContentAccessor{},
		PodcastAccessor:   &testPodcastAccessor{},
		AdvertisementAccessor: &testAdvertisementAccessor{
			isUserEligibleForAdvertisement: true,
			advertisements:                 []advertising.Advertisement{},
		},
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
	if body.Advertisement != nil {
		t.Errorf("Expected no advertisement, but there was one")
	}
}
