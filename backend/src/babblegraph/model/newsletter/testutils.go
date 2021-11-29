package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/util/math/decimal"
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"babblegraph/util/text"
	"babblegraph/wordsmith"
	"fmt"
	"strings"
)

func isIDExcluded(id documents.DocumentID, excludedIDs []documents.DocumentID) bool {
	for _, i := range excludedIDs {
		if i == id {
			return true
		}
	}
	return false
}

func isDomainValid(domain string, validDomains []string) bool {
	for _, d := range validDomains {
		if d == domain {
			return true
		}
	}
	return false
}

func containsTopic(topic contenttopics.ContentTopic, topics []contenttopics.ContentTopic) bool {
	for _, t := range topics {
		if t == topic {
			return true
		}
	}
	return false
}

func containsLemma(lemma wordsmith.LemmaID, description string) bool {
	tokens := text.TokenizeUnique(description)
	for _, t := range tokens {
		if t == lemma.Str() {
			return true
		}
	}
	return false
}

type getDefaultDocumentInput struct {
	Topics                 []contenttopics.ContentTopic
	Lemmas                 []wordsmith.LemmaID
	SeedJobIngestTimestamp *int64
}

func getDefaultDocumentWithLink(idx int, emailRecordID email.ID, userAccessor userPreferencesAccessor, input getDefaultDocumentInput) (*documents.DocumentWithScore, *Link, error) {
	var lemmatizedDescription *string
	if len(input.Lemmas) > 0 {
		var descriptionParts []string
		for _, lemma := range input.Lemmas {
			descriptionParts = append(descriptionParts, lemma.Str())
		}
		lemmatizedDescription = ptr.String(strings.Join(descriptionParts, " "))
	}
	doc := documents.Document{
		ID:               documents.DocumentID(fmt.Sprintf("web_doc-%d", idx)),
		Version:          documents.Version4,
		URL:              fmt.Sprintf("https://www.elmundo.es/%d", idx),
		ReadabilityScore: 50,
		LanguageCode:     wordsmith.LanguageCodeSpanish,
		DocumentType:     documents.TypeArticle,
		Metadata: documents.Metadata{
			Title:       ptr.String(fmt.Sprintf("Document %d", idx)),
			Image:       ptr.String(fmt.Sprintf("https://www.elmundo.es/%d.jpg", idx)),
			URL:         ptr.String(fmt.Sprintf("https://www.elmundo.es/%d", idx)),
			Description: ptr.String(fmt.Sprintf("This is document #%d", idx)),
		},
		Domain:                 "elmundo.es",
		Topics:                 input.Topics,
		HasPaywall:             ptr.Bool(false),
		LemmatizedDescription:  lemmatizedDescription,
		SeedJobIngestTimestamp: input.SeedJobIngestTimestamp,
	}
	link, err := makeLinkFromDocument(emailRecordID, userAccessor, doc)
	if err != nil {
		return nil, nil, err
	}
	return &documents.DocumentWithScore{
		Score:    decimal.FromInt64(1),
		Document: doc,
	}, link, nil
}

func testCategory(expected, result Category) error {
	var errs []string
	matchedLinks := make(map[documents.DocumentID]bool)
	for _, expectedLink := range expected.Links {
		var didFindLink bool
		for _, resultLink := range result.Links {
			isSameLink, err := testLink(expectedLink, resultLink)
			if isSameLink {
				matchedLinks[resultLink.DocumentID] = true
				if err != nil {
					errs = append(errs, err.Error())
				}
				didFindLink = true
				break
			}
		}
		if !didFindLink {
			errs = append(errs, fmt.Sprintf("Expected link for document ID %s, but didn't get it", expectedLink.DocumentID))
		}
	}
	for _, resultLink := range result.Links {
		if _, ok := matchedLinks[resultLink.DocumentID]; !ok {
			errs = append(errs, fmt.Sprintf("Got link for document ID %s, but didn't expect it", resultLink.DocumentID))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}

func testLink(expected, result Link) (bool, error) {
	var errs []string
	if expected.DocumentID != result.DocumentID {
		return false, nil
	}
	if err := testutils.CompareNullableString(expected.ImageURL, result.ImageURL); err != nil {
		errs = append(errs, fmt.Sprintf("Image URL for link %s: %s", expected.DocumentID, err.Error()))
	}
	if err := testutils.CompareNullableString(expected.Title, result.Title); err != nil {
		errs = append(errs, fmt.Sprintf("Title for link %s: %s", expected.DocumentID, err.Error()))
	}
	if err := testutils.CompareNullableString(expected.Description, result.Description); err != nil {
		errs = append(errs, fmt.Sprintf("Description for link %s: %s", expected.DocumentID, err.Error()))
	}
	if len(errs) > 0 {
		return true, fmt.Errorf(strings.Join(errs, "\n"))
	}
	return true, nil
}
