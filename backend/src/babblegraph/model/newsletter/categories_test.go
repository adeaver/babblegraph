package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"babblegraph/util/text"
	"babblegraph/wordsmith"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDefaultCategories(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	emailRecordID := email.NewEmailRecordID()
	documentTopics := []contenttopics.ContentTopic{
		contenttopics.ContentTopicArt,
		contenttopics.ContentTopicAstronomy,
		contenttopics.ContentTopicArchitecture,
		contenttopics.ContentTopicAutomotive,
		contenttopics.ContentTopicCulture,
	}
	userAccessor := &testUserAccessor{
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
			contenttopics.ContentTopicAstronomy,
			contenttopics.ContentTopicArchitecture,
			contenttopics.ContentTopicAutomotive,
		},
	}
	var expectedCategories []Category
	var docs []documents.DocumentWithScore
	for idx, topic := range documentTopics {
		doc, link, err := getDefaultDocumentWithLink(idx, emailRecordID, userAccessor, getDefaultDocumentInput{
			Topics: []contenttopics.ContentTopic{topic},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		if containsTopic(topic, userAccessor.getUserTopics()) {
			displayName, err := contenttopics.ContentTopicNameToDisplayName(topic)
			if err != nil {
				t.Fatalf("Error setting up test: %s", err.Error())
			}
			expectedCategories = append(expectedCategories, Category{
				Name: ptr.String(text.ToTitleCaseForLanguage(displayName.Str(), wordsmith.LanguageCodeSpanish)),
				Links: []Link{
					*link,
				},
			})
		}
		docs = append(docs, *doc)
	}
	categories, err := getDocumentCategories(c, getDocumentCategoriesInput{
		emailRecordID: emailRecordID,
		languageCode:  wordsmith.LanguageCodeSpanish,
		userAccessor:  userAccessor,
		docsAccessor: &testDocsAccessor{
			documents: docs,
		},
		numberOfDocumentsInNewsletter: ptr.Int(4),
	})
	if err != nil {
		t.Fatalf("Got error %s", err.Error())
	}
	if len(categories) != 4 {
		t.Errorf("Expected 4 categories, but got %d", len(categories))
	}
	var errs []string
	matchedCategories := make(map[string]bool)
	for _, e := range expectedCategories {
		var didFindCategory bool
		for _, c := range categories {
			switch {
			case c.Name == nil:
				errs = append(errs, "Got null category name, but did not expect one")
			case *c.Name == *e.Name:
				if err := testCategory(e, c); err != nil {
					errs = append(errs, fmt.Sprintf("Error on category %s: %s", *e.Name, err.Error()))
				}
				didFindCategory = true
				matchedCategories[*c.Name] = true
				break
			}
		}
		if !didFindCategory {
			errs = append(errs, fmt.Sprintf("Expected category %s, but didn't get it", *e.Name))
		}
	}
	for _, c := range categories {
		if c.Name != nil {
			if _, ok := matchedCategories[*c.Name]; !ok {
				errs = append(errs, fmt.Sprintf("Got category %s, but didn't expect it", *c.Name))
			}
		}
	}
	if len(errs) > 0 {
		t.Errorf(strings.Join(errs, "\n"))
	}
}

func TestGenericCategory(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	emailRecordID := email.NewEmailRecordID()
	documentTopics := []contenttopics.ContentTopic{
		contenttopics.ContentTopicArt,
		contenttopics.ContentTopicAstronomy,
		contenttopics.ContentTopicArchitecture,
		contenttopics.ContentTopicAutomotive,
		contenttopics.ContentTopicCulture,
	}
	userAccessor := &testUserAccessor{
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userTopics: []contenttopics.ContentTopic{},
	}
	var expectedLinks []Link
	var docs []documents.DocumentWithScore
	for idx, topic := range documentTopics {
		doc, link, err := getDefaultDocumentWithLink(idx, emailRecordID, userAccessor, getDefaultDocumentInput{
			Topics: []contenttopics.ContentTopic{topic},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		expectedLinks = append(expectedLinks, *link)
		docs = append(docs, *doc)
	}
	categories, err := getDocumentCategories(c, getDocumentCategoriesInput{
		emailRecordID: emailRecordID,
		languageCode:  wordsmith.LanguageCodeSpanish,
		userAccessor:  userAccessor,
		docsAccessor: &testDocsAccessor{
			documents: docs,
		},
		numberOfDocumentsInNewsletter: ptr.Int(4),
	})
	if err != nil {
		t.Fatalf("Got error %s", err.Error())
	}
	switch {
	case len(categories) != 1:
		t.Errorf("Expected 1 category, but got %d", len(categories))
	case categories[0].Name != nil:
		t.Errorf("Expected category to have null name, but got %s", *categories[0].Name)
	case len(categories[0].Links) != 4:
		t.Errorf("Expected category to have 4 links, but got %d", len(categories[0].Links))
	}
}

func TestCategoryWithGeneric(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	emailRecordID := email.NewEmailRecordID()
	documentTopics := []contenttopics.ContentTopic{
		contenttopics.ContentTopicArt,
		contenttopics.ContentTopicAstronomy,
		contenttopics.ContentTopicArchitecture,
		contenttopics.ContentTopicAutomotive,
		contenttopics.ContentTopicCulture,
	}
	userAccessor := &testUserAccessor{
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userTopics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
	}
	var expectedCategories []Category
	var docs []documents.DocumentWithScore
	for idx, topic := range documentTopics {
		doc, link, err := getDefaultDocumentWithLink(idx, emailRecordID, userAccessor, getDefaultDocumentInput{
			Topics: []contenttopics.ContentTopic{topic},
		})
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		if containsTopic(topic, userAccessor.getUserTopics()) {
			displayName, err := contenttopics.ContentTopicNameToDisplayName(topic)
			if err != nil {
				t.Fatalf("Error setting up test: %s", err.Error())
			}
			expectedCategories = append(expectedCategories, Category{
				Name: ptr.String(text.ToTitleCaseForLanguage(displayName.Str(), wordsmith.LanguageCodeSpanish)),
				Links: []Link{
					*link,
				},
			})
		}
		docs = append(docs, *doc)
	}
	categories, err := getDocumentCategories(c, getDocumentCategoriesInput{
		emailRecordID: emailRecordID,
		languageCode:  wordsmith.LanguageCodeSpanish,
		userAccessor:  userAccessor,
		docsAccessor: &testDocsAccessor{
			documents: docs,
		},
		numberOfDocumentsInNewsletter: ptr.Int(4),
	})
	if err != nil {
		t.Fatalf("Got error %s", err.Error())
	}
	if len(categories) != 2 {
		t.Errorf("Expected 2 category, but got %d", len(categories))
	}
	if err := testutils.CompareNullableString(categories[0].Name, expectedCategories[0].Name); err != nil {
		t.Errorf("Error on category name: %s", err.Error())
	}
	if len(categories[0].Links) != 1 {
		t.Errorf("Expected first category to have 1 link, but got %d", len(categories[0].Links))
	}
	genericCategoryDisplayName := ptr.String(text.ToTitleCaseForLanguage(contenttopics.GenericCategoryNameForLanguage(wordsmith.LanguageCodeSpanish).Str(), wordsmith.LanguageCodeSpanish))
	if err := testutils.CompareNullableString(categories[1].Name, genericCategoryDisplayName); err != nil {
		t.Errorf("Error on generic category name: %s", err.Error())
	}
	if len(categories[1].Links) != 3 {
		t.Errorf("Expected generic category to have 3 links, but got %d", len(categories[1].Links))
	}
}

func TestFavorRecentDocuments(t *testing.T) {
	c := ctx.GetDefaultLogContext()
	emailRecordID := email.NewEmailRecordID()
	userAccessor := &testUserAccessor{
		readingLevel: &userReadingLevel{
			LowerBound: 30,
			UpperBound: 80,
		},
		userTopics: []contenttopics.ContentTopic{},
	}
	var expectedLinks []Link
	var docs []documents.DocumentWithScore
	for idx := 0; idx <= 8; idx++ {
		doc, link, err := getDefaultDocumentWithLink(idx, emailRecordID, userAccessor, getDefaultDocumentInput{
			Topics: []contenttopics.ContentTopic{contenttopics.ContentTopicArt},
		})
		doc.Document.SeedJobIngestTimestamp = ptr.Int64(time.Now().Add(time.Duration(-2*(8-idx)*24) * time.Hour).Unix())
		if err != nil {
			t.Fatalf("Error setting up test: %s", err.Error())
		}
		expectedLinks = append(expectedLinks, *link)
		docs = append(docs, *doc)
	}
	categories, err := getDocumentCategories(c, getDocumentCategoriesInput{
		emailRecordID: emailRecordID,
		languageCode:  wordsmith.LanguageCodeSpanish,
		userAccessor:  userAccessor,
		docsAccessor: &testDocsAccessor{
			documents: docs,
		},
		numberOfDocumentsInNewsletter: ptr.Int(4),
	})
	if err != nil {
		t.Fatalf("Got error %s", err.Error())
	}
	switch {
	case len(categories) != 1:
		t.Errorf("Expected 1 category, but got %d", len(categories))
	case categories[0].Name != nil:
		t.Errorf("Expected category to have no name, but got %s", *categories[0].Name)
	}
	for _, link := range categories[0].Links {
		originalIdx, err := strconv.Atoi(strings.TrimPrefix(link.DocumentID.Str(), "web_doc-"))
		switch {
		case err != nil:
			t.Errorf("Got error converting string to index: %s", err.Error())
		case originalIdx <= 4:
			t.Errorf("Expected only recent documents (idx 5-8), but got document with idx %d", originalIdx)
		}
	}
}
