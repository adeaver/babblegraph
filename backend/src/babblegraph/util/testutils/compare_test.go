package testutils

import (
	"babblegraph/util/ptr"
	"fmt"
	"testing"
)

func TestCompareStringLists(t *testing.T) {
	type testCase struct {
		a   []string
		b   []string
		err error
	}

	testCases := []testCase{
		{
			a:   []string{"a", "b", "b"},
			b:   []string{"a", "b"},
			err: fmt.Errorf("First had [b], Second had []"),
		}, {
			a:   []string{"a", "b", "b", "b"},
			b:   []string{"a", "b"},
			err: fmt.Errorf("First had [b b], Second had []"),
		}, {
			a:   []string{"b", "a", "b"},
			b:   []string{"a", "b"},
			err: fmt.Errorf("First had [b], Second had []"),
		}, {
			a:   []string{"b", "a"},
			b:   []string{"a", "b"},
			err: nil,
		}, {
			a:   []string{"a", "b"},
			b:   []string{"a", "b"},
			err: nil,
		},
	}
	for idx, tc := range testCases {
		result := CompareStringLists(tc.a, tc.b)
		if result.Error() != tc.err.Error() {
			t.Errorf("Error on test case %d. Expected %s. Got %s", idx+1, tc.err.Error(), result.Error())
		}
	}
}

func TestCompareNullableStrings(t *testing.T) {
	type testCase struct {
		input    *string
		expected *string
		err      error
	}

	testCases := []testCase{
		{
			input:    ptr.String("hello"),
			expected: nil,
			err:      fmt.Errorf("Expected null, but got hello"),
		}, {
			input:    nil,
			expected: ptr.String("hello"),
			err:      fmt.Errorf("Expected hello, but got null"),
		}, {
			input:    ptr.String("goodbye"),
			expected: ptr.String("hello"),
			err:      fmt.Errorf("Expected hello, but got goodbye"),
		}, {
			input:    ptr.String("hello"),
			expected: ptr.String("hello"),
			err:      nil,
		},
	}
	for idx, tc := range testCases {
		result := CompareNullableString(tc.input, tc.expected)
		if result.Error() != tc.err.Error() {
			t.Errorf("Error on test case %d. Expected %s. Got %s", idx+1, tc.err.Error(), result.Error())
		}
	}
}

func TestCompareStringMapSimple(t *testing.T) {
	type testCase struct {
		a   map[string]string
		b   map[string]string
		err error
	}

	testCases := []testCase{
		{
			a:   map[string]string{"a": "1", "b": "2"},
			b:   map[string]string{"a": "1"},
			err: fmt.Errorf("b is missing key b, a has value 2"),
		}, {
			a:   map[string]string{"a": "1", "b": "2"},
			b:   map[string]string{"a": "1", "b": "3"},
			err: fmt.Errorf("mismatch on key b, b has value 3, a has value 2"),
		}, {
			a:   map[string]string{"a": "1"},
			b:   map[string]string{"a": "1", "b": "2"},
			err: fmt.Errorf("a is missing key b, b has value 2"),
		}, {
			a:   map[string]string{"a": "1", "b": "2"},
			b:   map[string]string{"a": "1", "b": "2"},
			err: nil,
		},
	}
	for idx, tc := range testCases {
		result := CompareStringMap(tc.a, tc.b)
		if result.Error() != tc.err.Error() {
			t.Errorf("Error on test case %d. Expected %s. Got %s", idx+1, tc.err.Error(), result.Error())
		}
	}
}
