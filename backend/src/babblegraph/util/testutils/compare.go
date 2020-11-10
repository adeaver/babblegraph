package testutils

import (
	"fmt"
	"strings"
)

func CompareStringLists(a, b []string) error {
	mapA := make(map[string]int)
	for _, s := range a {
		num, _ := mapA[s]
		mapA[s] = num + 1
	}
	var bDiff []string
	for _, s := range b {
		num, ok := mapA[s]
		switch {
		case !ok:
			bDiff = append(bDiff, s)
		case num == 1:
			delete(mapA, s)
		default:
			mapA[s] = num - 1
		}
	}
	var aDiff []string
	for s, num := range mapA {
		for i := 0; i < num; i++ {
			aDiff = append(aDiff, s)
		}
	}
	if len(aDiff) == 0 && len(bDiff) == 0 {
		return nil
	}
	return fmt.Errorf("First had %v, Second had %v", aDiff, bDiff)
}

func CompareNullableString(result, expected *string) error {
	switch {
	case result == nil && expected != nil:
		return fmt.Errorf("Expected %s, but got null", *expected)
	case result != nil && expected == nil:
		return fmt.Errorf("Expected null, but got %s", *result)
	case result == nil && expected == nil,
		*result == *expected:
		return nil
	case *result != *expected:
		return fmt.Errorf("Expected %s, but got %s", *expected, *result)
	default:
		panic("unreachable")
	}
}

func CompareStringMap(a, b map[string]string) error {
	var diff []string
	for key, value := range a {
		bValue, ok := b[key]
		switch {
		case !ok:
			diff = append(diff, fmt.Sprintf("b is missing key %s, a has value %v", key, value))
		case bValue != value:
			diff = append(diff, fmt.Sprintf("mismatch on key %s, b has value %v, a has value %v", key, bValue, value))
		}
	}
	for key, value := range b {
		_, ok := a[key]
		if !ok {
			diff = append(diff, fmt.Sprintf("a is missing key %s, b has value %v", key, value))
		}
	}
	if len(diff) > 0 {
		return fmt.Errorf(strings.Join(diff, "\n"))
	}
	return nil
}
