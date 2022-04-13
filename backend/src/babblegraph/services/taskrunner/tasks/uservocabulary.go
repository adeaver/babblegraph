package tasks

import (
	"babblegraph/model/userlemma"
	"babblegraph/model/users"
	"babblegraph/model/uservocabulary"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

func MigrateUserVocabulary(c ctx.LogContext) error {
	return database.WithTx(func(tx *sqlx.Tx) error {
		u, err := users.GetAllActiveUsers(tx)
		if err != nil {
			return err
		}
		for _, user := range u {
			mappings, err := userlemma.GetVisibleMappingsForUser(tx, user.ID)
			if err != nil {
				return err
			}
			spotlightRecords, err := userlemma.GetLemmaReinforcementRecordsForUserOrderedBySentOn(tx, user.ID)
			if err != nil {
				return err
			}
			lemmaIDToEntryID := make(map[wordsmith.LemmaID]uservocabulary.UserVocabularyEntryID)
			for _, m := range mappings {
				var lemma *wordsmith.Lemma
				if err := wordsmith.WithWordsmithTx(func(wordsmithTx *sqlx.Tx) error {
					var err error
					lemma, err = wordsmith.GetLemmaByID(wordsmithTx, m.LemmaID)
					return err
				}); err != nil {
					return err
				}
				h := &uservocabulary.HashableLemma{
					LemmaID:   lemma.ID,
					LemmaText: lemma.LemmaText,
				}
				userVocabularyEntryID, err := uservocabulary.UpsertVocabularyEntry(tx, uservocabulary.UpsertVocabularyEntryInput{
					UserID:       user.ID,
					LanguageCode: m.LanguageCode,
					IsActive:     m.IsActive,
					IsVisible:    true,
					Hashable:     h,
				})
				if err != nil {
					return err
				}
				lemmaIDToEntryID[lemma.ID] = *userVocabularyEntryID
			}
			for _, spotlight := range spotlightRecords {
				userVocabularyEntryID, ok := lemmaIDToEntryID[spotlight.LemmaID]
				if !ok {
					c.Infof("Lemma %s did not yield a user vocabulary entry", spotlight.LemmaID)
					continue
				}
				if err := uservocabulary.CreateUserVocabularyFromSpotlight(tx, userVocabularyEntryID, spotlight); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
