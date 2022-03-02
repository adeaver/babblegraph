package contentingestion

import (
	"babblegraph/model/content"
	"babblegraph/model/links2"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

func processWebsiteHTML1Link(c ctx.LogContext, link links2.Link) error {
	var shouldMarkAsComplete bool
	defer func() {
		if shouldMarkAsComplete {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return links2.SetURLAsFetched(tx, link.URLIdentifier)
			}); err != nil {
				c.Errorf("Error marking link %s complete: %s", link.URLIdentifier, err.Error())
			}
		}
	}()
	c.Infof("Processing new link with URL %s", link.URL)
	u := link.URL
	p := urlparser.ParseURL(u)
	if p == nil {
		c.Infof("Received link that did not parse")
		shouldMarkAsComplete = true
		return nil
	}
	var shouldExit bool
	var source *content.Source
	err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		shouldExit, err = content.IsParsedURLASeedURL(tx, *p)
		if err != nil {
			return err
		}
		source, err = content.LookupActiveSourceForSourceID(c, tx, *link.SourceID)
		return err
	})
	switch {
	case err != nil:
		c.Warnf("Error verifying if url %s is a seed url: %s", u, err.Error())
		return err
	case shouldExit == true:
		shouldMarkAsComplete = true
		return nil
	case source == nil:
		c.Infof("No source found")
		return nil
	default:
		// no-op
	}

}
