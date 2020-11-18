package userreadability

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const lookupUserReadabilityForLanguageQuery = "SELECT * FROM user_readability_level WHERE user_id = $1 AND language_code = $2"

func lookupUserReadabilityForLanguage(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) ([]userReadabilityLevel, error) {
	var matches []userReadabilityLevel
	if err := tx.Select(&matches, lookupUserReadabilityForLanguageQuery, userID, languageCode); err != nil {
		return nil, err
	}
	return matches, nil
}