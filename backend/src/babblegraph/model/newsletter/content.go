package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type contentAccessor interface {
	GetSourceByID(sourceID content.SourceID) (*content.Source, error)
	GetDisplayNameByTopicID(topicID content.TopicID) (*string, error)
}

type DefaultContentAccessor struct {
	allowableSourcesByID       map[content.SourceID]content.Source
	topicDisplayNamesByTopicID map[content.TopicID]string
}

func GetDefaultContentAccessor(tx *sqlx.Tx, languageCode wordsmith.LanguageCode) (*DefaultContentAccessor, error) {
	sources, err := content.GetAllowableSources(tx)
	if err != nil {
		return nil, err
	}
	sourcesByID := make(map[content.SourceID]content.Source)
	for _, source := range sources {
		sourcesByID[source.ID] = source
	}
	topicDisplayNames, err := content.GetTopicDisplayNamesForLanguage(tx, languageCode)
	if err != nil {
		return nil, err
	}
	displayNamesByID := make(map[content.TopicID]string)
	for _, d := range topicDisplayNames {
		displayNamesByID[d.TopicID] = d.Label
	}
	return &DefaultContentAccessor{
		topicDisplayNamesByTopicID: displayNamesByID,
		allowableSourcesByID:       sourcesByID,
	}, nil
}

func (d *DefaultContentAccessor) GetSourceByID(sourceID content.SourceID) (*content.Source, error) {
	source, ok := d.allowableSourcesByID[sourceID]
	if !ok {
		return nil, fmt.Errorf("Source %s not found", sourceID)
	}
	return &source, nil
}

func (d *DefaultContentAccessor) GetDisplayNameByTopicID(topicID content.TopicID) (*string, error) {
	displayName, ok := d.topicDisplayNamesByTopicID[topicID]
	if !ok {
		return nil, fmt.Errorf("No display name found for topic id %s", topicID)
	}
	return ptr.String(displayName), nil
}
