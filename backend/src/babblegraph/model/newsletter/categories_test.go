package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/util/ptr"
	"babblegraph/util/text"
	"babblegraph/wordsmith"
	"fmt"
	"strings"
	"testing"
)

func TestDefaultCategories(t *testing.T) {
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
		docs = append(docs, *doc)
	}
	categories, err := getDocumentCategories(getDocumentCategoriesInput{
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
