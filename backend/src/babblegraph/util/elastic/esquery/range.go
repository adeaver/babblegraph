package esquery

type rangeQueryBuilder struct {
	fieldName            string
	GreaterThan          *int64 `json:"gt,omitempty"`
	GreaterThanOrEqualTo *int64 `json:"gte,omitempty"`
	LessThan             *int64 `json:"lt,omitempty"`
	LessThanOrEqualTo    *int64 `json:"lte,omitempty"`
}

func NewRangeQueryBuilderForFieldName(fieldName string) *rangeQueryBuilder {
	return &rangeQueryBuilder{fieldName: fieldName}
}

func (r *rangeQueryBuilder) GreaterThanInt64(value int64) {
	r.GreaterThan = &value
}

func (r *rangeQueryBuilder) GreaterThanOrEqualToInt64(value int64) {
	r.GreaterThanOrEqualTo = &value
}

func (r *rangeQueryBuilder) LessThanInt64(value int64) {
	r.LessThan = &value
}

func (r *rangeQueryBuilder) LessThanOrEqualToInt64(value int64) {
	r.LessThanOrEqualTo = &value
}

func (r *rangeQueryBuilder) BuildRangeQuery() query {
	return makeQuery(queryNameRange.Str(), makeQuery(r.fieldName, r))
}
