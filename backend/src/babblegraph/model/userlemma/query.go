package userlemma

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const (
	getVisibleMappingsForUser       = "SELECT * FROM user_lemma_mappings WHERE user_id = $1 AND is_visible = TRUE"
	setOrToggleMappingsForUserQuery = "INSERT INTO user_lemma_mappings (user_id, lemma_id, language_code) VALUES ($1, $2, $3) ON CONFLICT (user_id, lemma_id) DO UPDATE SET is_visible = $4, is_active = $5"
	toggleMappingActiveStateQuery   = "UPDATE user_lemma_mappings SET is_active = $1 WHERE user_id = $2 AND lemma_id = $3"

	getLemmaReinforcementSpolightRecordForUserQuery = "SELECT * FROM user_lemma_reinforcement_spotlight_records WHERE user_id = $1 ORDER BY last_sent_on DESC"
	setLemmaReinforcementSpolightRecordForUserQuery = `INSERT INTO
        user_lemma_reinforcement_spotlight_records (user_id, language_code, lemma_id, last_sent_on, number_of_times_sent)
        VALUES ($1, $2, $3, timezone('utc', now()), $4)
    ON CONFLICT (user_id, lemma_id)
    DO UPDATE SET
        language_code=$2,
        last_sent_on=timezone('utc', now()),
        number_of_times_sent=$4`
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
	res, err := tx.Exec(setOrToggleMappingsForUserQuery, userID, lemmaID, languageCode, true, true)
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

func RemoveMappingForUser(tx *sqlx.Tx, userID users.UserID, lemmaID wordsmith.LemmaID, languageCode wordsmith.LanguageCode) (bool, error) {
	res, err := tx.Exec(setOrToggleMappingsForUserQuery, userID, lemmaID, languageCode, false, false)
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

func GetLemmaReinforcementRecordsForUserOrderedBySentOn(tx *sqlx.Tx, userID users.UserID) ([]UserLemmaReinforcementSpotlightRecord, error) {
	var matches []dbUserLemmaReinforcementSpotlightRecord
	if err := tx.Select(&matches, getLemmaReinforcementSpotlightRecordForUserQuery, userID); err != nil {
		return nil, err
	}
	var out []UserLemmaReinforcementSpotlightRecord
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

type UpsertLemmaReinforcementSpotlightRecordInput struct {
	UserID            users.UserID
	LemmaID           wordsmith.LemmaID
	LanguageCode      wordsmith.LanguageCode
	NumberOfTimesSent int64
}

func UpsertLemmaReinforcementSpotlightRecord(tx *sqlx.Tx, input UpsertLemmaReinforcementSpotlightRecordInput) error {
	if _, err := tx.Exec(setLemmaReinforcementSpotlightRecordForUserQuery, input.UserID, input.LanguageCode, input.LemmaID, input.NumberOfTimesSent); err != nil {
		return err
	}
	return nil
}
