package esquery

type sort map[string]interface{}

type sortBody struct {
	Order   sortOrder        `json:"order"`
	Missing sortMissingValue `json:"missing"`
}

type orderedSort struct {
	sorts []sort
}

func NewOrderedSort(sorts ...sort) *orderedSort {
	return &orderedSort{
		sorts: sorts,
	}
}

type sortOrder string

const (
	sortOrderAscending  sortOrder = "asc"
	sortOrderDescending sortOrder = "desc"
)

type sortMissingValue string

const (
	sortMissingValueFirst sortMissingValue = "_first"
	sortMissingValueLast  sortMissingValue = "_last"
)

type sortBuilder struct {
	fieldName string
	order     sortOrder
	missing   sortMissingValue
}

func NewAscendingSortBuilder(fieldName string) *sortBuilder {
	return &sortBuilder{
		fieldName: fieldName,
		order:     sortOrderAscending,
		missing:   sortMissingValueLast,
	}
}

func NewDescendingSortBuilder(fieldName string) *sortBuilder {
	return &sortBuilder{
		fieldName: fieldName,
		order:     sortOrderDescending,
		missing:   sortMissingValueLast,
	}
}

func (s *sortBuilder) WithMissingValuesLast() {
	s.missing = sortMissingValueLast
}

func (s *sortBuilder) WithMissingValuesFirst() {
	s.missing = sortMissingValueFirst
}

func (s *sortBuilder) AsSort() sort {
	return sort(map[string]interface{}{
		s.fieldName: sortBody{
			Order:   s.order,
			Missing: s.missing,
		},
	})
}
