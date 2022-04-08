package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/util/geo"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"fmt"
	"strings"
)

const testSourceID content.SourceID = content.SourceID("test-source")

var testSource = content.Source{
	ID:                    testSourceID,
	Title:                 "Test Source",
	URL:                   "babblegraph.com",
	Type:                  content.SourceTypeNewsWebsite,
	Country:               geo.CountryCodeUnitedStates,
	IngestStrategy:        content.IngestStrategyWebsiteHTML1,
	LanguageCode:          wordsmith.LanguageCodeSpanish,
	IsActive:              true,
	ShouldUseURLAsSeedURL: true,
	MonthlyAccessLimit:    nil,
}

type testContentAccessor struct {
	topicIDsNotInList []content.TopicID
}

func (t *testContentAccessor) GetSourceByID(sourceID content.SourceID) (*content.Source, error) {
	if sourceID == testSourceID {
		return &testSource, nil
	}
	return nil, fmt.Errorf("Unsupported source ID %s", sourceID)
}

func (t *testContentAccessor) GetDisplayNameByTopicID(topicID content.TopicID) (*string, error) {
	parts := strings.Split(topicID.Str(), "-")
	return ptr.String(parts[len(parts)-1]), nil
}

func (t *testContentAccessor) GetAllTopicsNotInList(topicIDs []content.TopicID) []content.TopicID {
	return t.topicIDsNotInList
}
