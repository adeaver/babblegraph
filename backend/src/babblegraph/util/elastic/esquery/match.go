package esquery

func Match(key string, value interface{}) query {
	subquery := makeQuery(key, value)
	return makeQuery(queryNameMatch.Str(), subquery)
}
