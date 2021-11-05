package tasks

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/links2"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func BackfillContentTopicsLength() error {
	return database.WithTx(func(tx *sqlx.Tx) error {
		var countErrs int64
		err := links2.GetCursorForFetchVersion(tx, links2.FetchVersion7, func(link links2.Link) (bool, error) {
			return false, database.WithTx(func(tx *sqlx.Tx) error {
				log.Println(fmt.Sprintf("Processing document with URL %s", link.URL))
				contentTopics, err := contenttopics.GetTopicsForURL(tx, link.URL)
				if err != nil {
					log.Println(fmt.Sprintf("Document with URL %s has error getting content topics %s, skipping", link.URL, err.Error()))
					countErrs++
					return err
				}
				u := urlparser.ParseURL(link.URL)
				if u == nil {
					log.Println(fmt.Sprintf("Document with URL %s has invalid URL, skipping", link.URL))
					countErrs++
					return nil
				}
				if err := documents.UpdateDocumentForURL(*u, documents.UpdateDocumentInput{
					Version:      documents.Version7,
					TopicsLength: ptr.Int64(int64(len(contentTopics))),
				}); err != nil {
					log.Println(fmt.Sprintf("Document with URL %s has error updating %s, skipping", link.URL, err.Error()))
					countErrs++
					return err
				}
				return nil
			})
		})
		log.Println(fmt.Sprintf("Finished reindexing with %d errors", countErrs))
		return err
	})
}
