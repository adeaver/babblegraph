package indexing

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"log"
)

type IndexDocumentInput struct {
	ParsedHTMLPage  ingesthtml.ParsedHTMLPage
	TextMetadata    textprocessing.TextMetadata
	LanguageCode    wordsmith.LanguageCode
	DocumentVersion documents.Version
	URL             urlparser.ParsedURL
	TopicsForURL    []contenttopics.ContentTopic
}

func IndexDocument(input IndexDocumentInput) error {
	var lemmatizedDescriptionText *string
	var lemmatizedDescriptionIndexMappings []int
	if input.TextMetadata.LemmatizedDescription != nil {
		lemmatizedDescriptionText = ptr.String(input.TextMetadata.LemmatizedDescription.LemmatizedText)
		lemmatizedDescriptionIndexMappings = input.TextMetadata.LemmatizedDescription.IndexMappings
	}
	docID, err := documents.AssignIDAndIndexDocument(documents.IndexDocumentInput{
		URL:                                input.URL,
		ReadabilityScore:                   input.TextMetadata.ReadabilityScore.ToInt64Rounded(),
		LanguageCode:                       input.LanguageCode,
		Metadata:                           input.ParsedHTMLPage.Metadata,
		Topics:                             input.TopicsForURL,
		LemmatizedDescription:              lemmatizedDescriptionText,
		LemmatizedDescriptionIndexMappings: lemmatizedDescriptionIndexMappings,

		// These will get changed later
		Version: input.DocumentVersion,
		Type:    documents.TypeArticle,
	})
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Indexed url %s with ID: %s", input.URL, string(*docID)))
	return nil
}
