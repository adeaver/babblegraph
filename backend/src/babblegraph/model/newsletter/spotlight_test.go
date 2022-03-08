package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/userlemma"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSpotlightRecordsForUserWithAccount(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	expectedLemma := wordsmith.LemmaID("word3")
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
		userNewsletterPreferences: &usernewsletterpreferences.UserNewsletterPreferences{
			ShouldIncludeLemmaReinforcementSpotlight: true,
			LanguageCode:                             wordsmith.LanguageCodeSpanish,
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
				LemmaID:      expectedLemma,
				LastSentOn:   time.Now().Add(-8 * 24 * time.Hour),
			},
		},
		trackingLemmas: []wordsmith.LemmaID{
			"word1", "word2", expectedLemma,
		},
	}
	var correctLink *Link
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
		doc, link, err := getDefaultDocumentWithLink(c, i, emailRecordID, &testContentAccessor{}, userAccessor, getDefaultDocumentInput{
			Topics:                 []contenttopics.ContentTopic{contenttopics.ContentTopicArt},
			SeedJobIngestTimestamp: ptr.Int64(time.Now().Add(-1 * time.Duration(15-i) * 24 * time.Hour).Unix()),
			Lemmas:                 []wordsmith.LemmaID{lemma},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		docs = append(docs, *doc)
		if lemma.Str() == expectedLemma.Str() {
			correctLink = link
		}
	}
	wordsmithAccessor := &testWordsmithAccessor{
		lemmasByID: lemmasByID,
	}
	docsAccessor := &testDocsAccessor{documents: docs}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
	})
	switch {
	case err != nil:
		t.Fatalf(err.Error())
	case testNewsletter == nil:
		t.Errorf("Expected non-null newsletter, but it was not")
	case testNewsletter.Body.LemmaReinforcementSpotlight == nil:
		t.Errorf("Expected non-null newsletter lemma reinforcement, but it was not")
	default:
		if testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText != expectedLemma.Str() {
			t.Errorf("Expected lemma to be %s, but got %s", expectedLemma.Str(), testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText)
		}
		correctDocument, err := testLink(testNewsletter.Body.LemmaReinforcementSpotlight.Document, *correctLink)
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

func TestSpotlightRecordsForUserWithoutAccount(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	expectedLemma := wordsmith.LemmaID("word3")
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		userID:              users.UserID("abc123"),
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: false,
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userNewsletterPreferences: &usernewsletterpreferences.UserNewsletterPreferences{
			ShouldIncludeLemmaReinforcementSpotlight: true,
			LanguageCode:                             wordsmith.LanguageCodeSpanish,
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
				LemmaID:      expectedLemma,
				LastSentOn:   time.Now().Add(-8 * 24 * time.Hour),
			},
		},
		trackingLemmas: []wordsmith.LemmaID{
			"word1", "word2", expectedLemma,
		},
	}
	var correctLink *Link
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
		doc, link, err := getDefaultDocumentWithLink(c, i, emailRecordID, &testContentAccessor{}, userAccessor, getDefaultDocumentInput{
			Topics:                 []contenttopics.ContentTopic{contenttopics.ContentTopicArt},
			SeedJobIngestTimestamp: ptr.Int64(time.Now().Add(-1 * time.Duration(15-i) * 24 * time.Hour).Unix()),
			Lemmas:                 []wordsmith.LemmaID{lemma},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		docs = append(docs, *doc)
		if lemma.Str() == expectedLemma.Str() {
			correctLink = link
		}
	}
	wordsmithAccessor := &testWordsmithAccessor{
		lemmasByID: lemmasByID,
	}
	docsAccessor := &testDocsAccessor{documents: docs}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
	})
	switch {
	case err != nil:
		t.Fatalf(err.Error())
	case testNewsletter == nil:
		t.Errorf("Expected non-null newsletter, but it was not")
	case testNewsletter.Body.LemmaReinforcementSpotlight == nil:
		t.Errorf("Expected non-null newsletter lemma reinforcement, but it was not")
	default:
		if testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText != expectedLemma.Str() {
			t.Errorf("Expected lemma to be %s, but got %s", expectedLemma.Str(), testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText)
		}
		correctDocument, err := testLink(testNewsletter.Body.LemmaReinforcementSpotlight.Document, *correctLink)
		if !correctDocument {
			t.Errorf("Document ID from links is not correct")
		}
		if err != nil {
			t.Errorf("Error comparing links: %s", err.Error())
		}
		if !strings.Contains(testNewsletter.Body.LemmaReinforcementSpotlight.PreferencesLink, "manage") {
			t.Errorf("Expected link to contain manage, but got %s", testNewsletter.Body.LemmaReinforcementSpotlight.PreferencesLink)
		}
		if !strings.HasSuffix(testNewsletter.Body.LemmaReinforcementSpotlight.PreferencesLink, "preferences") {
			t.Errorf("Expected link to end with preferences, but got %s", testNewsletter.Body.LemmaReinforcementSpotlight.PreferencesLink)
		}
	}
}

func TestSpotlightRecordsForTrackedLemmaWithoutSpotlight(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	expectedLemma := wordsmith.LemmaID("word4")
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		userID:              users.UserID("abc123"),
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: false,
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userNewsletterPreferences: &usernewsletterpreferences.UserNewsletterPreferences{
			ShouldIncludeLemmaReinforcementSpotlight: true,
			LanguageCode:                             wordsmith.LanguageCodeSpanish,
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
		trackingLemmas: []wordsmith.LemmaID{
			"word1", "word2", "word3", expectedLemma,
		},
	}
	var correctLink *Link
	emailRecordID := email.NewEmailRecordID()
	lemmasByID := make(map[wordsmith.LemmaID]wordsmith.Lemma)
	var docs []documents.DocumentWithScore
	for i := 16; i >= 0; i-- {
		lemma := wordsmith.LemmaID(fmt.Sprintf("word%d", i))
		lemmasByID[lemma] = wordsmith.Lemma{
			ID:        lemma,
			Language:  wordsmith.LanguageCodeSpanish,
			LemmaText: lemma.Str(),
		}
		doc, link, err := getDefaultDocumentWithLink(c, i, emailRecordID, &testContentAccessor{}, userAccessor, getDefaultDocumentInput{
			Topics:                 []contenttopics.ContentTopic{contenttopics.ContentTopicArt},
			SeedJobIngestTimestamp: ptr.Int64(time.Now().Add(-1 * time.Duration(15-i) * 24 * time.Hour).Unix()),
			Lemmas:                 []wordsmith.LemmaID{lemma},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		docs = append(docs, *doc)
		if lemma.Str() == expectedLemma.Str() {
			correctLink = link
		}
	}
	wordsmithAccessor := &testWordsmithAccessor{
		lemmasByID: lemmasByID,
	}
	docsAccessor := &testDocsAccessor{documents: docs}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor: wordsmithAccessor,
		EmailAccessor:     emailAccessor,
		UserAccessor:      userAccessor,
		DocsAccessor:      docsAccessor,
	})
	switch {
	case err != nil:
		t.Fatalf(err.Error())
	case testNewsletter == nil:
		t.Errorf("Expected non-null newsletter, but it was not")
	case testNewsletter.Body.LemmaReinforcementSpotlight == nil:
		t.Errorf("Expected non-null newsletter lemma reinforcement, but it was not")
	default:
		if testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText != expectedLemma.Str() {
			t.Errorf("Expected lemma to be %s, but got %s", expectedLemma.Str(), testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText)
		}
		correctDocument, err := testLink(testNewsletter.Body.LemmaReinforcementSpotlight.Document, *correctLink)
		if !correctDocument {
			t.Errorf("Document ID from links is not correct")
		}
		if err != nil {
			t.Errorf("Error comparing links: %s", err.Error())
		}
	}
}
