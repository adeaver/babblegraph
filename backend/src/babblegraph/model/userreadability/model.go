package userreadability

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
)

type userReadabilityVersion int

const (
	version1 userReadabilityVersion = 1
)

type userReadabilityLevelID string

type userReadabilityLevel struct {
	ID               userReadabilityLevelID `db:"_id"`
	LanguageCode     wordsmith.LanguageCode `db:"language_code"`
	UserID           users.UserID           `db:"user_id"`
	ReadabilityLevel int64                  `db:"readability_level"`
	Version          userReadabilityVersion `db:"version"`
}
