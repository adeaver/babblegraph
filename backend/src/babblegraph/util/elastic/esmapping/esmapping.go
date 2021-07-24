package esmapping

import (
	"babblegraph/util/elastic"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
)

type Mapping map[string]mappingBody

type mappingBody struct {
	MappingOptions
	Type mappingType `json:"type"`
}

// TODO: add the following options:
// - index_options
// - index_prefixes
// - similarity
// - store
// - term_vector
type MappingOptions struct {
	Analyzer       *string `json:"analyzer,omitempty"`
	CopyTo         *string `json:"copy_to,omitempty"`
	Dynamic        *bool   `json:"dynamic,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty"`
	Format         *Format `json:"format,omitempty"`
	IgnoreAbove    *int64  `json:"ignore_above,omitempty"`
	Index          *bool   `json:"index,omitempty"`
	Meta           *Meta   `json:"meta,omitempty"`
	Normalizer     *string `json:"normalizer,omitempty"`
	NullValue      *string `json:"null_value,omitempty"`
	SearchAnalyzer *string `json:"search_analyzer,omitempty"`

	// This should be created by a mapping object
	Properties *Mapping `json:"properties,omitempty"`

	// This should be created using the MappingWithFields function
	Fields *Mapping `json:"fields,omitempty"`

	// The following fields are possible options, but generally
	// should be avoided. Any code using these fields should be
	// thoroughly commented.
	Boost                *int64 `json:"boost,omitempty"`
	Coerce               *bool  `json:"coerce,omitempty"`
	DocValues            *bool  `json:"doc_values,omitempty"`
	EagerGlobalOrdinals  *bool  `json:"eager_global_ordinals,omitempty"`
	FieldData            *bool  `json:"field_data,omitempty"`
	IgnoreMalformed      *bool  `json:"ignore_malformed,omitempty"`
	IndexPhrases         *bool  `json:"index_phrases,omitempty"`
	Norms                *bool  `json:"norms,omitempty"`
	PositionIncrementGap *int64 `json:"position_increment_gap,omitempty"`
}

type mappingType string

const (
	mappingTypeObject  mappingType = "object"
	mappingTypeText    mappingType = "text"
	mappingTypeBoolean mappingType = "boolean"
	mappingTypeKeyword mappingType = "keyword"
	mappingTypeLong    mappingType = "long"
	mappingTypeDate    mappingType = "date"
)

func (m mappingType) Ptr() *mappingType {
	return &m
}

// https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-date-format.html
type Format string

type Meta struct {
	Unit       *Unit       `json:"unit,omitempty"`
	MetricType *MetricType `json:"metric_type,omitempty"`
}

type Unit string

type MetricType string

func makeMapping(fieldName string, mType mappingType, options MappingOptions) Mapping {
	body := mappingBody{
		MappingOptions: options,
		Type:           mType,
	}
	return Mapping(map[string]mappingBody{
		fieldName: body,
	})
}

func MappingWithFields(m Mapping, fields []Mapping) Mapping {
	flattenedFields, err := flattenMappings(fields)
	if err != nil {
		panic(fmt.Errorf("Error making mapping with fields for mapping %+v: %s", m, err.Error()))
	}
	for fieldName, mappingBody := range m {
		options := mappingBody.MappingOptions
		options.Fields = flattenedFields
		return makeMapping(fieldName, mappingBody.Type, options)
	}
	panic(fmt.Errorf("Unreachable statement while creating mapping with fields"))
}

func flattenMappings(mappings []Mapping) (*Mapping, error) {
	properties := make(map[string]mappingBody)
	for _, m := range mappings {
		for fieldName, mBody := range m {
			if _, ok := properties[fieldName]; ok {
				return nil, fmt.Errorf("duplicate field name %s in mapping request", fieldName)
			}
			properties[fieldName] = mBody
		}
	}
	m := Mapping(properties)
	return &m, nil
}

type updateMappingsRequestBody struct {
	Properties map[string]mappingBody `json:"properties"`
}

func UpdateMapping(index elastic.Index, mappings []Mapping) error {
	properties, err := flattenMappings(mappings)
	if err != nil {
		return err
	}
	body := updateMappingsRequestBody{
		Properties: *properties,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return elastic.RunUpdateMappingsRequest(esapi.IndicesPutMappingRequest{
		Index: []string{index.GetName()},
		Body:  strings.NewReader(string(bodyBytes)),
	}, esapi.IndicesPutMappingRequest{
		Index: []string{index.GetName()},
		Body:  strings.NewReader(string(bodyBytes)),
	})
}
