package links

import "github.com/jmoiron/sqlx"

func GetUnfetchedLinks(tx *sqlx.Tx) (map[Domain][]Link, error) {
	var unfetchedLinks []dbLink
	if err := tx.Select(&unfetchedLinks, "SELECT * FROM links WHERE has_fetched=FALSE ORDER BY position ASC"); err != nil {
		return nil, err
	}
	out := make(map[Domain][]Link)
	for _, link := range unfetchedLinks {
		linksForDomain, ok := out[link.Domain]
		if !ok {
			linksForDomain = []Link{}
		}
		linksForDomain = append(linksForDomain, link.ToLink())
		out[link.Domain] = linksForDomain
	}
	return out, nil
}

const insertLinkQuery = "INSERT INTO links (domain, url) VALUES ($1, $2) ON CONFLICT DO NOTHING"

func InsertLinks(tx *sqlx.Tx, links []Link) error {
	for _, l := range links {
		if _, err := tx.Exec(insertLinkQuery, l.Domain, l.URL); err != nil {
			return err
		}
	}
	return nil
}

func SetURLAsFetched(tx *sqlx.Tx, url string) error {
	if _, err := tx.Exec("UPDATE links SET has_fetched=TRUE WHERE url=$1", url); err != nil {
		return err
	}
	return nil
}
