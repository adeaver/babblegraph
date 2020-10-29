package links

type dbLink struct {
	ID         LinkID `db:"_id"`
	Domain     Domain `db:"domain"`
	URL        string `db:"url"`
	HasFetched bool   `db:"has_fetched"`
	Position   int    `db:"position"`
}

func (d dbLink) ToLink() Link {
	return Link{
		Domain: d.Domain,
		URL:    d.URL,
	}
}

type Link struct {
	Domain Domain
	URL    string
}

type LinkID string

type Domain string

func (d Domain) Str() string {
	return string(d)
}
