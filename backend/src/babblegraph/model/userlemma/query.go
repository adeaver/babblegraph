package userlemma

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const (
	getVisibleMappingsForUser     = "SELECT * FROM user_lemma_mappings WHERE user_id = $1 AND is_visible = TRUE"
	addMappingsForUserQuery       = "INSERT INTO user_lemma_mappings (user_id, lemma_id, language_code) VALUES ($1, $2, $3) ON CONFLICT (user_id, lemma_id) DO UPDATE SET is_visible = TRUE, is_active = TRUE"
	toggleMappingActiveStateQuery = "UPDATE user_lemma_mappings SET is_active = $1 WHERE user_id = $2 AND lemma_id = $3"
)

func GetVisibleMappingsForUser(tx *sqlx.Tx, userID users.UserID) ([]Mapping, error) {
	var matches []dbMapping
	if err := tx.Select(&matches, getVisibleMappingsForUser, userID); err != nil {
		return nil, err
	}
	var out []Mapping
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func AddMappingForUser(tx *sqlx.Tx, userID users.UserID, lemmaID wordsmith.LemmaID, languageCode wordsmith.LanguageCode) (bool, error) {
	res, err := tx.Exec(addMappingsForUserQuery, userID, lemmaID, languageCode)
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

func ToggleMappingActiveState(tx *sqlx.Tx, userID users.UserID, lemmaID wordsmith.LemmaID, currentState bool) (bool, error) {
	nextState := "TRUE"
	if currentState {
		nextState = "FALSE"
	}
	res, err := tx.Exec(toggleMappingActiveStateQuery, nextState, userID, lemmaID)
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
