package esquery

func IDs(ids []string) query {
	subquery := makeQuery("values", ids)
	return makeQuery(queryNameIDs.Str(), subquery)
}
