package esmapping

func MakeTextMapping(fieldName string, options MappingOptions) Mapping {
	return makeMapping(fieldName, mappingTypeText, options)
}

func MakeKeywordMapping(fieldName string, options MappingOptions) Mapping {
	return makeMapping(fieldName, mappingTypeKeyword, options)
}

func MakeBooleanMapping(fieldName string, options MappingOptions) Mapping {
	return makeMapping(fieldName, mappingTypeBoolean, options)
}

func MakeLongMapping(fieldName string, options MappingOptions) Mapping {
	return makeMapping(fieldName, mappingTypeLong, options)
}

func MakeDateMapping(fieldName string, options MappingOptions) Mapping {
	return makeMapping(fieldName, mappingTypeDate, options)
}
