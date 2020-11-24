package bufferedfetch

import (
	"fmt"
	"testing"
)

func TestSingleCall(t *testing.T) {
	numCalls := 0
	if err := Register("test-single", func() (interface{}, error) {
		numCalls++
		var ints []int
		for i := 0; i < 5; i++ {
			ints = append(ints, i)
		}
		return ints, nil
	}); err != nil {
		t.Errorf("Got error: %s", err.Error())
	}
	for i := 0; i < 5; i++ {
		var v int
		if err := WithNextBufferedValue("test-single", func(i interface{}) error {
			var ok bool
			v, ok = i.(int)
			if !ok {
				return fmt.Errorf("Incorrect type")
			}
			return nil
		}); err != nil {
			t.Errorf("Got error: %s", err.Error())
		}
		if v != i {
			t.Errorf("should have gotten %d, but got %d", v, i)
		}
	}
	if numCalls != 1 {
		t.Errorf("should have only called refill function 1 time, but called %d times", numCalls)
	}
}

func TestDoubleCall(t *testing.T) {
	numCalls := 0
	if err := Register("test-double", func() (interface{}, error) {
		numCalls++
		var ints []int
		for i := 0; i < 5; i++ {
			ints = append(ints, i)
		}
		return ints, nil
	}); err != nil {
		t.Errorf("Got error: %s", err.Error())
	}
	for i := 0; i < 6; i++ {
		var v int
		if err := WithNextBufferedValue("test-double", func(i interface{}) error {
			var ok bool
			v, ok = i.(int)
			if !ok {
				return fmt.Errorf("Incorrect type")
			}
			return nil
		}); err != nil {
			t.Errorf("Got error: %s", err.Error())
		}
		expected := i
		if i == 5 {
			expected = 0
		}
		if v != expected {
			t.Errorf("should have gotten %d, but got %d", v, expected)
		}
	}
	if numCalls != 2 {
		t.Errorf("should have only called refill function 2 times, but called %d times", numCalls)
	}
}
