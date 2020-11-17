package userreadability

import (
	"babblegraph/model/userreadability/spanishreadability"
	"babblegraph/model/users"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
	"fmt"
)

type GetReadabilityScoreRangeForUser struct {
	UserID       users.UserID
	LanguageCode wordsmith.LanguageCode
}

type ReadabilityScoreRange struct {
	MinScore decimal.Number
	MaxScore decimal.Number
}

func GetReadabilityScoreRangeForUser(tx *sqlx.Tx, input GetReadabilityScoreRangeForUser) (*ReadabilityScoreRange, error) {
	readabilityLevels, err := lookupUserReadabilityForLanguage(tx, input.UserID, input.LanguageCode)
	if err != nil {
		return nil, err
	}
	var readabilityLevel *userReadabilityLevel
	for _, level := range readabilityLevels {
		if level.Version == version1 {
			readabilityLevel = &level
			// It's important to break here
			// So that readabilityLevel doesn't get clobbered
			// since level is a pointer
			break
		}
	}
	if readabilityLevel == nil {
		return nil, fmt.Errorf("No readability level found")
	}
	switch input.LanguageCode {
	case wordsmith.LanguageCodeSpanish:
		levelAsInt := int(readabilityLevel.ReadabilityLevel.ToInt64())
		minScore, maxScore, err := spanishreadability.GetReadabilityScoreRangeForLevel(levelAsInt)
		if err != nil {
			return nil, err
		}
		return &ReadabilityScoreRange{
			MinScore: *minScore,
			MaxScore: *maxScore,
		}, nil
	default:
		panic(fmt.Sprintf("Unrecognized language %s", input.LanguageCode.Str()))
	}
}
