package esquery

func Terms(field string, terms []string) query {
	subquery := makeQuery(field, terms)
	return makeQuery(queryNameTerms.Str(), subquery)
}
