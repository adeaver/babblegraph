package esquery

func MatchAll() query {
	return makeQuery("match_all", struct{}{})
}
