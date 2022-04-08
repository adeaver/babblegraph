package int2

import (
	"babblegraph/util/ptr"
	"fmt"
)

func MinInt(i ...int) (*int, error) {
	if len(i) == 0 {
		return nil, fmt.Errorf("Empty list of ints")
	}
	minInt := i[0]
	for idx := 1; idx < len(i); idx++ {
		j := i[idx]
		if j < minInt {
			minInt = j
		}
	}
	return ptr.Int(minInt), nil
}

func MustMinInt(i ...int) int {
	minInt, err := MinInt(i...)
	if err != nil {
		panic(fmt.Sprintf("Error getting max int: %s", err.Error()))
	}
	return *minInt
}

func MaxInt(i ...int) (*int, error) {
	if len(i) == 0 {
		return nil, fmt.Errorf("Empty list of ints")
	}
	maxInt := i[0]
	for idx := 1; idx < len(i); idx++ {
		j := i[idx]
		if j > maxInt {
			maxInt = j
		}
	}
	return ptr.Int(maxInt), nil
}

func MustMaxInt(i ...int) int {
	maxInt, err := MaxInt(i...)
	if err != nil {
		panic(fmt.Sprintf("Error getting max int: %s", err.Error()))
	}
	return *maxInt
}
