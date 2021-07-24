package esmapping

import "fmt"

func MakeObjectMapping(objectName string, properties []Mapping) Mapping {
	objectProperties := make(map[string]mappingBody)
	for _, m := range properties {
		for fieldName, body := range m {
			if _, ok := objectProperties[fieldName]; ok {
				panic(fmt.Errorf("Object with field name %s has duplicate mapping defined for %s", objectName, fieldName))
			}
			objectProperties[fieldName] = body
		}
	}
	asMapping := Mapping(objectProperties)
	return makeMapping(objectName, mappingTypeObject, MappingOptions{
		Properties: &asMapping,
	})
}
