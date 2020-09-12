package wordsmith

import "github.com/jmoiron/sqlx"

func GetWords(words []string, language LanguageCode) ([]Word, error) {
	var matches []dbWord
	if err := withTx(func(tx *sqlx.Tx) error {
		query, args, err := sqlx.In("SELECT * FROM words WHERE language=? AND word IN (?)", language.Str(), words)
		if err != nil {
			return err
		}
		sql := tx.Rebind(query)
		return tx.Select(&matches, sql, args...)
	}); err != nil {
		return nil, err
	}
	var out []Word
	for _, match := range matches {
		out = append(out, match.ToWord())
	}
	return out, nil
}
