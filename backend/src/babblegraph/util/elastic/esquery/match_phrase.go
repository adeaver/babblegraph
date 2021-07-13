package esquery

func MatchPhrase(key string, value interface{}) query {
	subquery := makeQuery(key, value)
	return makeQuery(queryNameMatchPhrase.Str(), subquery)
}
