package spanishreadability

import (
	"babblegraph/util/math/decimal"
	"fmt"
)

type readabilityScoreRange struct {
	minScore decimal.Number
	maxScore decimal.Number
}

var readabilityLevelsToScore = map[int]readabilityScoreRange{
	// Elementary School Level, very easy
	1: readabilityScoreRange{
		minScore: decimal.FromInt64(90),
		maxScore: decimal.FromInt64(100),
	},
	// 6th Grade, easy to read, conversational
	2: readabilityScoreRange{
		minScore: decimal.FromInt64(80),
		maxScore: decimal.FromInt64(90),
	},
	// 7th Grade, fairly easy to read
	3: readabilityScoreRange{
		minScore: decimal.FromInt64(70),
		maxScore: decimal.FromInt64(80),
	},
	// 8th & 9th Grade, plain spanish, understood by high schoolers
	4: readabilityScoreRange{
		minScore: decimal.FromInt64(60),
		maxScore: decimal.FromInt64(70),
	},
	// 10th to 12th Grade, fairly difficult, pre-college
	5: readabilityScoreRange{
		minScore: decimal.FromInt64(60),
		maxScore: decimal.FromInt64(70),
	},
	// College level
	6: readabilityScoreRange{
		minScore: decimal.FromInt64(30),
		maxScore: decimal.FromInt64(50),
	},
	// College-graduate level
	7: readabilityScoreRange{
		minScore: decimal.FromInt64(10),
		maxScore: decimal.FromInt64(30),
	},
	// Professional level
	8: readabilityScoreRange{
		minScore: decimal.FromInt64(0),
		maxScore: decimal.FromInt64(10),
	},
}

func GetReadabilityScoreRangeForLevel(level int) (_min, _max *decimal.Number, _err error) {
	levelRange, ok := readabilityLevelsToScore[level]
	if !ok {
		return nil, nil, fmt.Errorf("Invalid level for Spanish %d", level)
	}
	return &levelRange.minScore, &levelRange.maxScore, nil
}
