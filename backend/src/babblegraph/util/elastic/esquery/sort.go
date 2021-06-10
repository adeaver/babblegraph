package esquery

type sort map[string]interface{}

type sortBody struct {
	Order        sortOrder         `json:"order"`
	Missing      *sortMissingValue `json:"missing,omitempty"`
	UnmappedType *sortUnmappedType `json:"unmapped_type,omitempty"`
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

func (s sortMissingValue) Ptr() *sortMissingValue {
	return &s
}

type sortUnmappedType string

const (
	sortUnmappedTypeLong sortUnmappedType = "long"
)

func (s sortUnmappedType) Ptr() *sortUnmappedType {
	return &s
}

type sortBuilder struct {
	fieldName    string
	order        sortOrder
	missing      *sortMissingValue
	unmappedType *sortUnmappedType
}

func NewAscendingSortBuilder(fieldName string) *sortBuilder {
	return &sortBuilder{
		fieldName: fieldName,
		order:     sortOrderAscending,
	}
}

func NewDescendingSortBuilder(fieldName string) *sortBuilder {
	return &sortBuilder{
		fieldName: fieldName,
		order:     sortOrderDescending,
	}
}

func (s *sortBuilder) WithMissingValuesLast() {
	s.missing = sortMissingValueLast.Ptr()
}

func (s *sortBuilder) WithMissingValuesFirst() {
	s.missing = sortMissingValueFirst.Ptr()
}

func (s *sortBuilder) AsUnmappedTypeLong() {
	s.unmappedType = sortUnmappedTypeLong.Ptr()
}

func (s *sortBuilder) AsSort() sort {
	return sort(map[string]interface{}{
		s.fieldName: sortBody{
			Order:        s.order,
			Missing:      s.missing,
			UnmappedType: s.unmappedType,
		},
	})
}
