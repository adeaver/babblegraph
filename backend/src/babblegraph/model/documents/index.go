package documents

import (
	"babblegraph/util/deref"
	"babblegraph/util/elastic"
	"babblegraph/util/opengraph"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"time"
)

const documentIndexName string = "web_documents"

type documentIndex struct{}

func (d documentIndex) GetName() string {
	return documentIndexName
}

func (d documentIndex) ValidateDocument(document interface{}) error {
	if _, ok := document.(Document); !ok {
		return fmt.Errorf("could not validate interface %+v, to be of type web_document", document)
	}
	return nil
}

func (d documentIndex) GenerateIDForDocument(document interface{}) (*string, error) {
	doc, ok := document.(Document)
	if !ok {
		return nil, fmt.Errorf("could not validate interface %+v, to be of type web_document", document)
	}
	docID := makeDocumentIndexForURL(urlparser.MustParseURL(doc.URL))
	return ptr.String(string(docID)), nil
}

type IndexDocumentInput struct {
	URL              urlparser.ParsedURL
	Metadata         map[string]string
	Type             Type
	Version          Version
	LanguageCode     wordsmith.LanguageCode
	LemmatizedBody   string
	ReadabilityScore int64
}

func AssignIDAndIndexDocument(input IndexDocumentInput) (*DocumentID, error) {
	documentID := makeDocumentIndexForURL(input.URL)
	ogMetadata := opengraph.GetBasicMetadata(input.Metadata)
	if err := elastic.IndexDocument(documentIndex{}, Document{
		ID:               documentID,
		Version:          input.Version,
		URL:              input.URL.URL,
		ReadabilityScore: input.ReadabilityScore,
		LemmatizedBody:   input.LemmatizedBody,
		LanguageCode:     input.LanguageCode,
		DocumentType:     input.Type.Ptr(),
		Domain:           input.URL.Domain,
		Metadata: &Metadata{
			Title:              deref.String(ogMetadata.Title, ""),
			Image:              deref.String(ogMetadata.ImageURL, ""),
			URL:                deref.String(ogMetadata.URL, ""),
			Description:        deref.String(ogMetadata.Description, ""),
			PublicationTimeUTC: getPublicationTimeUTCOrNil(ogMetadata.PublicationTime),
		},
	}); err != nil {
		return nil, err
	}
	return &documentID, nil
}

func getPublicationTimeUTCOrNil(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	return ptr.Time(t.UTC())
}
