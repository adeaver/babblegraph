package esmapping

func MakeTextMapping(fieldName string, options MappingOptions) Mapping {
	return makeMapping(fieldName, mappingTypeText, options)
}
