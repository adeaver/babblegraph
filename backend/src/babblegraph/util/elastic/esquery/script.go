package esquery

func Script(script string) query {
	subquery := makeQuery(queryNameScript.Str(), script)
	return makeQuery(queryNameScript.Str(), subquery)
}
