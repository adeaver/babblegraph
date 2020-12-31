package wordsmith

import "github.com/jmoiron/sqlx"

func GetSortedRankingsForWords(languageCode LanguageCode, words []string) ([]WordRanking, error) {
	var matches []dbWordRanking
	if err := withTx(func(tx *sqlx.Tx) error {
		corpusID := getCurrentCorpusIDForLanguageCode(languageCode)
		query, args, err := sqlx.In("SELECT * FROM word_rankings WHERE corpus_id=? AND word IN (?) ORDER BY ranking ASC", corpusID, words)
		if err != nil {
			return err
		}
		sql := tx.Rebind(query)
		return tx.Select(&matches, sql, args...)
	}); err != nil {
		return nil, err
	}
	var out []WordRanking
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}
