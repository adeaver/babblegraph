package indexing

import (
	"babblegraph/model/documents"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/wordsmith"
	"fmt"
	"log"
)

type IndexDocumentInput struct {
	ParsedHTMLPage  ingesthtml.ParsedHTMLPage
	TextMetadata    textprocessing.TextMetadata
	LanguageCode    wordsmith.LanguageCode
	DocumentVersion documents.Version
	URL             string
}

func IndexDocument(input IndexDocumentInput) error {
	documentMetadata := documents.ExtractMetadataFromMap(input.ParsedHTMLPage.Metadata)
	docID, err := documents.AssignIDAndIndexDocument(&documents.Document{
		URL:              input.URL,
		Version:          input.DocumentVersion,
		ReadabilityScore: input.TextMetadata.ReadabilityScore.ToInt64Rounded(),
		LanguageCode:     input.LanguageCode,
		LemmatizedBody:   input.TextMetadata.LemmatizedText,
		DocumentType:     documents.TypeArticle.Ptr(),
		Metadata:         &documentMetadata,
	})
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Indexed url %s with ID: %s", input.URL, string(*docID)))
	return nil
}
