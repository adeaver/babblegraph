package userreadability

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const (
	lookupUserReadabilityForLanguageQuery  = "SELECT * FROM user_readability_level WHERE user_id = $1 AND language_code = $2"
	lookupAllUserReadabilitiesForUserQuery = "SELECT * FROM user_readability_level WHERE user_id = $1 AND version = $2"
	updateUserReadability                  = "UPDATE user_readability_level SET readability_level = $1 WHERE user_id = $2 AND language_code = $3 AND version = $4"
)

func lookupUserReadabilityForLanguage(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) ([]userReadabilityLevel, error) {
	var matches []userReadabilityLevel
	if err := tx.Select(&matches, lookupUserReadabilityForLanguageQuery, userID, languageCode); err != nil {
		return nil, err
	}
	return matches, nil
}

func lookupUserReadabilitiesForUser(tx *sqlx.Tx, userID users.UserID) ([]userReadabilityLevel, error) {
	var matches []userReadabilityLevel
	if err := tx.Select(&matches, lookupAllUserReadabilitiesForUserQuery, userID, version1); err != nil {
		return nil, err
	}
	return matches, nil
}

func updateUserReadabilityForUser(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode, level int) (bool, error) {
	res, err := tx.Exec(updateUserReadability, level, userID, languageCode, version1)
	if err != nil {
		return false, err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	didUpdate := numRows > 0
	return didUpdate, nil
}
