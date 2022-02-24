package indexing

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
)

type IndexDocumentInput struct {
	ParsedHTMLPage         ingesthtml.ParsedHTMLPage
	TextMetadata           textprocessing.TextMetadata
	LanguageCode           wordsmith.LanguageCode
	DocumentVersion        documents.Version
	URL                    urlparser.ParsedURL
	SourceID               *content.SourceID
	TopicsForURL           []contenttopics.ContentTopic
	TopicIDs               []content.TopicID
	TopicMappingIDs        []content.TopicMappingID
	SeedJobIngestTimestamp *int64
}

func IndexDocument(c ctx.LogContext, input IndexDocumentInput) error {
	var lemmatizedDescriptionText *string
	var lemmatizedDescriptionIndexMappings []int
	if input.TextMetadata.LemmatizedDescription != nil {
		lemmatizedDescriptionText = ptr.String(input.TextMetadata.LemmatizedDescription.LemmatizedText)
		lemmatizedDescriptionIndexMappings = input.TextMetadata.LemmatizedDescription.IndexMappings
	}
	docID, err := documents.AssignIDAndIndexDocument(c, documents.IndexDocumentInput{
		URL:                                input.URL,
		SourceID:                           input.SourceID,
		ReadabilityScore:                   input.TextMetadata.ReadabilityScore.ToInt64Rounded(),
		LanguageCode:                       input.LanguageCode,
		Metadata:                           input.ParsedHTMLPage.Metadata,
		Topics:                             input.TopicsForURL,
		TopicIDs:                           input.TopicIDs,
		TopicMappingIDs:                    input.TopicMappingIDs,
		LemmatizedDescription:              lemmatizedDescriptionText,
		LemmatizedDescriptionIndexMappings: lemmatizedDescriptionIndexMappings,
		SeedJobIngestTimestamp:             input.SeedJobIngestTimestamp,
		HasPaywall:                         input.ParsedHTMLPage.IsPaywalled,

		// These will get changed later
		Version: input.DocumentVersion,
		Type:    documents.TypeArticle,
	})
	if err != nil {
		return err
	}
	c.Infof("Indexed url %s with ID: %s", input.URL, string(*docID))
	return nil
}
