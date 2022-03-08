package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/util/geo"
	"babblegraph/wordsmith"
	"fmt"
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

type testContentAccessor struct{}

func (t *testContentAccessor) GetSourceByID(sourceID content.SourceID) (*content.Source, error) {
	if sourceID == testSourceID {
		return &testSource, nil
	}
	return nil, fmt.Errorf("Unsupported source ID %s", sourceID)
}
