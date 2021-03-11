package esquery

func Term(key string, value interface{}) query {
	subquery := makeQuery(key, value)
	return makeQuery(queryNameTerm.Str(), subquery)
}
