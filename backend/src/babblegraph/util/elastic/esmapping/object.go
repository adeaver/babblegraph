package esmapping

import "fmt"

func MakeObjectMapping(objectName string, properties []Mapping) Mapping {
	objectProperties, err := flattenMappings(properties)
	if err != nil {
		panic(fmt.Errorf("Error making object mapping %s: %s", objectName, err.Error()))
	}
	return makeMapping(objectName, mappingTypeObject, MappingOptions{
		Properties: objectProperties,
	})
}
