package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/model/uservocabulary"
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
	expectedEntryID := uservocabulary.UserVocabularyEntryID("word3")
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
		userNewsletterPreferences: &usernewsletterpreferences.UserNewsletterPreferences{
			ShouldIncludeLemmaReinforcementSpotlight: true,
			LanguageCode:                             wordsmith.LanguageCodeSpanish,
		},
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
		allowableSourceIDs: []content.SourceID{
			content.SourceID("test-source"),
		},
		spotlightRecords: []uservocabulary.UserVocabularySpotlightRecord{
			{
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: "word1",
				LastSentOn:        time.Now(),
			}, {
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: "word2",
				LastSentOn:        time.Now(),
			}, {
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: expectedEntryID,
				LastSentOn:        time.Now().Add(-8 * 24 * time.Hour),
			},
		},
		vocabularyEntries: []uservocabulary.UserVocabularyEntry{
			{
				ID:                "word1",
				VocabularyID:      ptr.String("word1"),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				VocabularyDisplay: "word1",
			}, {
				ID:                "word2",
				VocabularyID:      ptr.String("word2"),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				VocabularyDisplay: "word2",
			}, {
				VocabularyID:      ptr.String(string(expectedEntryID)),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				ID:                expectedEntryID,
				VocabularyDisplay: string(expectedEntryID),
			},
		},
	}
	var correctLink *Link
	emailRecordID := email.NewEmailRecordID()
	contentAccessor := &testContentAccessor{}
	var docs []documents.DocumentWithScore
	for i := 15; i >= 0; i-- {
		lemma := wordsmith.LemmaID(fmt.Sprintf("word%d", i))
		doc, link, err := getDefaultDocumentWithLink(c, i, emailRecordID, contentAccessor, userAccessor, getDefaultDocumentInput{
			Topics: []content.TopicID{
				content.TopicID("topicid-art"),
			},
			SeedJobIngestTimestamp: ptr.Int64(time.Now().Add(-1 * time.Duration(15-i) * 24 * time.Hour).Unix()),
			Lemmas:                 []wordsmith.LemmaID{lemma},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		docs = append(docs, *doc)
		if lemma.Str() == expectedEntryID.Str() {
			correctLink = link
		}
	}
	wordsmithAccessor := &testWordsmithAccessor{}
	docsAccessor := &testDocsAccessor{documents: docs}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor:     wordsmithAccessor,
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          docsAccessor,
		ContentAccessor:       contentAccessor,
		AdvertisementAccessor: &testAdvertisementAccessor{},
		PodcastAccessor:       &testPodcastAccessor{},
	})
	switch {
	case err != nil:
		t.Fatalf(err.Error())
	case testNewsletter == nil:
		t.Errorf("Expected non-null newsletter, but it was not")
	case testNewsletter.Body.LemmaReinforcementSpotlight == nil:
		t.Errorf("Expected non-null newsletter lemma reinforcement, but it was not")
	default:
		if testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText != expectedEntryID.Str() {
			t.Errorf("Expected lemma to be %s, but got %s", expectedEntryID.Str(), testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText)
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
	expectedEntryID := uservocabulary.UserVocabularyEntryID("word3")
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		userID:              users.UserID("abc123"),
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: false,
		userTopics: []content.TopicID{
			content.TopicID("topicid-art"),
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userNewsletterPreferences: &usernewsletterpreferences.UserNewsletterPreferences{
			ShouldIncludeLemmaReinforcementSpotlight: true,
			LanguageCode:                             wordsmith.LanguageCodeSpanish,
		},
		allowableSourceIDs: []content.SourceID{
			content.SourceID("test-source"),
		},
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
		spotlightRecords: []uservocabulary.UserVocabularySpotlightRecord{
			{
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: "word1",
				LastSentOn:        time.Now(),
			}, {
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: "word2",
				LastSentOn:        time.Now(),
			}, {
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: expectedEntryID,
				LastSentOn:        time.Now().Add(-8 * 24 * time.Hour),
			},
		},
		vocabularyEntries: []uservocabulary.UserVocabularyEntry{
			{
				ID:                "word1",
				VocabularyID:      ptr.String("word1"),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				VocabularyDisplay: "word1",
			}, {
				ID:                "word2",
				VocabularyID:      ptr.String("word2"),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				VocabularyDisplay: "word2",
			}, {
				VocabularyID:      ptr.String(string(expectedEntryID)),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				ID:                expectedEntryID,
				VocabularyDisplay: string(expectedEntryID),
			},
		},
	}
	var correctLink *Link
	emailRecordID := email.NewEmailRecordID()
	contentAccessor := &testContentAccessor{}
	var docs []documents.DocumentWithScore
	for i := 15; i >= 0; i-- {
		lemma := wordsmith.LemmaID(fmt.Sprintf("word%d", i))
		doc, link, err := getDefaultDocumentWithLink(c, i, emailRecordID, contentAccessor, userAccessor, getDefaultDocumentInput{
			Topics: []content.TopicID{
				content.TopicID("topicid-art"),
			},
			SeedJobIngestTimestamp: ptr.Int64(time.Now().Add(-1 * time.Duration(15-i) * 24 * time.Hour).Unix()),
			Lemmas:                 []wordsmith.LemmaID{lemma},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		docs = append(docs, *doc)
		if lemma.Str() == expectedEntryID.Str() {
			correctLink = link
		}
	}
	wordsmithAccessor := &testWordsmithAccessor{}
	docsAccessor := &testDocsAccessor{documents: docs}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor:     wordsmithAccessor,
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          docsAccessor,
		ContentAccessor:       contentAccessor,
		AdvertisementAccessor: &testAdvertisementAccessor{},
		PodcastAccessor:       &testPodcastAccessor{},
	})
	switch {
	case err != nil:
		t.Fatalf(err.Error())
	case testNewsletter == nil:
		t.Errorf("Expected non-null newsletter, but it was not")
	case testNewsletter.Body.LemmaReinforcementSpotlight == nil:
		t.Errorf("Expected non-null newsletter lemma reinforcement, but it was not")
	default:
		if testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText != expectedEntryID.Str() {
			t.Errorf("Expected lemma to be %s, but got %s", expectedEntryID.Str(), testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText)
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
	expectedEntryID := uservocabulary.UserVocabularyEntryID("word4")
	emailAccessor := getTestEmailAccessor()
	userAccessor := &testUserAccessor{
		userID:              users.UserID("abc123"),
		languageCode:        wordsmith.LanguageCodeSpanish,
		doesUserHaveAccount: false,
		userTopics: []content.TopicID{
			content.TopicID("topicid-art"),
		},
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userNewsletterPreferences: &usernewsletterpreferences.UserNewsletterPreferences{
			ShouldIncludeLemmaReinforcementSpotlight: true,
			LanguageCode:                             wordsmith.LanguageCodeSpanish,
		},
		userNewsletterSchedule: usernewsletterpreferences.TestNewsletterSchedule{
			SendRequested:     true,
			NumberOfDocuments: 4,
		},
		allowableSourceIDs: []content.SourceID{
			content.SourceID("test-source"),
		},
		spotlightRecords: []uservocabulary.UserVocabularySpotlightRecord{
			{
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: "word1",
				LastSentOn:        time.Now(),
			}, {
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: "word2",
				LastSentOn:        time.Now(),
			}, {
				LanguageCode:      wordsmith.LanguageCodeSpanish,
				VocabularyEntryID: "word3",
				LastSentOn:        time.Now().Add(-8 * 24 * time.Hour),
			},
		},
		vocabularyEntries: []uservocabulary.UserVocabularyEntry{
			{
				ID:                "word1",
				VocabularyID:      ptr.String("word1"),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				VocabularyDisplay: "word1",
			}, {
				ID:                "word2",
				VocabularyID:      ptr.String("word2"),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				VocabularyDisplay: "word2",
			}, {
				VocabularyID:      ptr.String(string(expectedEntryID)),
				VocabularyType:    uservocabulary.VocabularyTypeLemma,
				ID:                expectedEntryID,
				VocabularyDisplay: string(expectedEntryID),
			},
		},
	}
	var correctLink *Link
	emailRecordID := email.NewEmailRecordID()
	contentAccessor := &testContentAccessor{}
	var docs []documents.DocumentWithScore
	for i := 16; i >= 0; i-- {
		lemma := wordsmith.LemmaID(fmt.Sprintf("word%d", i))
		doc, link, err := getDefaultDocumentWithLink(c, i, emailRecordID, contentAccessor, userAccessor, getDefaultDocumentInput{
			Topics: []content.TopicID{
				content.TopicID("topicid-art"),
			},
			SeedJobIngestTimestamp: ptr.Int64(time.Now().Add(-1 * time.Duration(15-i) * 24 * time.Hour).Unix()),
			Lemmas:                 []wordsmith.LemmaID{lemma},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		docs = append(docs, *doc)
		if lemma.Str() == expectedEntryID.Str() {
			correctLink = link
		}
	}
	wordsmithAccessor := &testWordsmithAccessor{}
	docsAccessor := &testDocsAccessor{documents: docs}
	testNewsletter, err := CreateNewsletter(c, CreateNewsletterInput{
		WordsmithAccessor:     wordsmithAccessor,
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          docsAccessor,
		ContentAccessor:       contentAccessor,
		AdvertisementAccessor: &testAdvertisementAccessor{},
		PodcastAccessor:       &testPodcastAccessor{},
	})
	switch {
	case err != nil:
		t.Fatalf(err.Error())
	case testNewsletter == nil:
		t.Errorf("Expected non-null newsletter, but it was not")
	case testNewsletter.Body.LemmaReinforcementSpotlight == nil:
		t.Errorf("Expected non-null newsletter lemma reinforcement, but it was not")
	default:
		if testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText != expectedEntryID.Str() {
			t.Errorf("Expected lemma to be %s, but got %s", expectedEntryID.Str(), testNewsletter.Body.LemmaReinforcementSpotlight.LemmaText)
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
