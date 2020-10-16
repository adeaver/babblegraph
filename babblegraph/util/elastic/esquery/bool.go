package esquery

type boolQueryBuilder struct {
	Must    []query `json:"must,omitempty"`
	MustNot []query `json:"must_not,omitempty"`
	Should  []query `json:"should,omitempty"`
	Filters []query `json:"filters,omitempty"`
}

func NewBoolQueryBuilder() *boolQueryBuilder {
	return &boolQueryBuilder{}
}

func (b *boolQueryBuilder) AddMust(q query) {
	b.Must = append(b.Must, q)
}

func (b *boolQueryBuilder) AddMustNot(q query) {
	b.MustNot = append(b.MustNot, q)
}

func (b *boolQueryBuilder) AddShould(q query) {
	b.Should = append(b.Should, q)
}

func (b *boolQueryBuilder) AddFilter(q query) {
	b.Filters = append(b.Filters, q)
}

func (b *boolQueryBuilder) BuildBoolQuery() query {
	return makeQuery(queryNameBool.Str(), b)
}