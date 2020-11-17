package userreadability

import (
	"babblegraph/model/users"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
)

type userReadabilityVersion int

const (
	version1 userReadabilityVersion = 1
)

type userReadabilityLevelID string

type userReadabilityLevel struct {
	ID               UserReadabilityLevelID `json:"id"`
	LanguageCode     wordsmith.LanguageCode `json:"language_code"`
	UserID           users.UserID           `json:"user_id"`
	ReadabilityLevel decimal.Number         `json:"readability_level"`
	Version          userReadabilityVersion `json:"version"`
}
