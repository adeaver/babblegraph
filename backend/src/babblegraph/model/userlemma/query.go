package userlemma

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const (
	getActiveMappingsForUserQuery = "SELECT * FROM user_lemma_mappings WHERE user_id = $1 AND is_active = TRUE"
	addMappingsForUserQuery       = "INSERT INTO user_lemma_mappings (user_id, lemma_id, language_code) VALUES ($1, $2, $3) ON CONFLICT DO UPDATE SET is_visible = TRUE, is_active = TRUE"
	setMappingAsInactiveQuery     = "UPDATE user_lemma_mappings SET is_active = FALSE WHERE user_id = $1 AND _id = $2"
	setMappingAsNotVisibleQuery   = "UPDATE user_lemma_mappings SET is_visible = FALSE WHERE user_id = $1 AND _id = $2"
)

func GetActiveMappingsForUser(tx *sqlx.Tx, userID users.UserID) ([]Mapping, error) {
	var matches []dbMapping
	if err := tx.Select(&matches, getActiveMappingsForUserQuery, userID); err != nil {
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

func SetMappingAsInactive(tx *sqlx.Tx, userID users.UserID, id MappingID) (bool, error) {
	res, err := tx.Exec(setMappingAsInactiveQuery, userID, id)
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

func SetMappingAsNotVisible(tx *sqlx.Tx, userID users.UserID, id MappingID) (bool, error) {
	res, err := tx.Exec(setMappingAsNotVisibleQuery, userID, id)
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
