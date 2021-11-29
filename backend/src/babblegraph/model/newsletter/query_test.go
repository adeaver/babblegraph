package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/userlemma"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"babblegraph/wordsmith"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestUserHasAccount(t *testing.T) {
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
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode: wordsmith.LanguageCodeSpanish,
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
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
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
	}
	docsAccessor := &testDocsAccessor{}
	testNewsletter, err := CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
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
	wordsmithAccessor := &testWordsmithAccessor{}
	emailAccessor := getTestEmailAccessor()
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
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
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
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
		userScheduleForDay: &usernewsletterschedule.UserNewsletterScheduleDayMetadata{
			IsActive: false,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
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

func TestSpotlightRecordsForUserWithAccount(t *testing.T) {
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: true,
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		spotlightRecords: []userlemma.UserLemmaReinforcementSpotlightRecord{
			{
				LanguageCode: wordsmith.LanguageCodeSpanish,
				LemmaID:      "word1",
				LastSentOn:   time.Now(),
			}, {
				LanguageCode: wordsmith.LanguageCodeSpanish,
				LemmaID:      "word2",
				LastSentOn:   time.Now(),
			}, {
				LanguageCode: wordsmith.LanguageCodeSpanish,
				LemmaID:      "word3",
				LastSentOn:   time.Now().Add(-8 * 24 * time.Hour),
			},
		},
	}
	var links []Link
	emailRecordID := email.NewEmailRecordID()
	lemmasByID := make(map[wordsmith.LemmaID]wordsmith.Lemma)
	var docs []documents.DocumentWithScore
	for i := 15; i >= 0; i-- {
		lemma := wordsmith.LemmaID(fmt.Sprintf("word%d", i))
		lemmasByID[lemma] = wordsmith.Lemma{
			ID:        lemma,
			Language:  wordsmith.LanguageCodeSpanish,
			LemmaText: lemma.Str(),
		}
		doc, link, err := getDefaultDocumentWithLink(i, emailRecordID, userAccessor, getDefaultDocumentInput{
			Topics:                 []contenttopics.ContentTopic{contenttopics.ContentTopicArt},
			SeedJobIngestTimestamp: ptr.Int64(time.Now().Add(-1 * time.Duration(15-i) * 24 * time.Hour).Unix()),
			Lemmas:                 []wordsmith.LemmaID{lemma},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		docs = append(docs, *doc)
		links = append(links, *link)
	}
	wordsmithAccessor := &testWordsmithAccessor{
		lemmasByID: lemmasByID,
	}
	docsAccessor := &testDocsAccessor{documents: docs}
	testNewsletter, err := CreateNewsletter(wordsmithAccessor, emailAccessor, userAccessor, docsAccessor)
	switch {
	case err != nil:
		t.Fatalf(err.Error())
	case testNewsletter == nil:
		t.Errorf("Expected non-null newsletter, but it was not")
	case testNewsletter.Body.LemmaReinforcementSpotlight == nil:
		t.Errorf("Expected non-null newsletter lemma reinforcement, but it was not")
	default:
		if testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText != "word3" {
			t.Errorf("Expected lemma to be word3, but got %s", testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText)
		}
		correctDocument, err := testLink(testNewsletter.Body.LemmaReinforcementSpotlight.Document, links[4])
		if !correctDocument {
			t.Errorf("Document ID from links is not correct")
		}
		if err != nil {
			t.Errorf("Error comparing links: %s", err.Error())
		}
		if !strings.HasSuffix(testNewsletter.Body.LemmaReinforcementSpotlight.PreferencesLink, "login?d=npf") {
			t.Errorf("Expected link to end with login?d=npf, but got %s", testNewsletter.Body.LemmaReinforcementSpotlight.PreferencesLink)
		}
	}
}
