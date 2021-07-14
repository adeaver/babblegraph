package esquery

type BoolQueryBuilder struct {
	Must    []query `json:"must,omitempty"`
	MustNot []query `json:"must_not,omitempty"`
	Should  []query `json:"should,omitempty"`
	Filter  []query `json:"filter,omitempty"`
}

func NewBoolQueryBuilder() *BoolQueryBuilder {
	return &BoolQueryBuilder{}
}

func (b *BoolQueryBuilder) AddMust(q query) {
	b.Must = append(b.Must, q)
}

func (b *BoolQueryBuilder) AddMustNot(q query) {
	b.MustNot = append(b.MustNot, q)
}

func (b *BoolQueryBuilder) AddShould(q query) {
	b.Should = append(b.Should, q)
}

func (b *BoolQueryBuilder) AddFilter(q query) {
	b.Filter = append(b.Filter, q)
}

func (b *BoolQueryBuilder) BuildBoolQuery() query {
	return makeQuery(queryNameBool.Str(), b)
}
