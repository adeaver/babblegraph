package deref

import "testing"

func TestDerefString(t *testing.T) {
	type testCase struct {
		input          *string
		expectedOutput string
	}
	testString := "hello"
	testCases := []testCase{
		{
			input:          nil,
			expectedOutput: "default string",
		}, {
			input:          &testString,
			expectedOutput: "hello",
		},
	}
	for idx, tc := range testCases {
		result := String(tc.input, "default string")
		if result != tc.expectedOutput {
			t.Errorf("Error on test case %d, expected %s, but got %s", idx+1, tc.expectedOutput, result)
		}
	}
}

func TestDerefInt(t *testing.T) {
	type testCase struct {
		input          *int
		expectedOutput int
	}
	testInt := 5
	testCases := []testCase{
		{
			input:          nil,
			expectedOutput: 42,
		}, {
			input:          &testInt,
			expectedOutput: 5,
		},
	}
	for idx, tc := range testCases {
		result := Int(tc.input, 42)
		if result != tc.expectedOutput {
			t.Errorf("Error on test case %d, expected %d, but got %d", idx+1, tc.expectedOutput, result)
		}
	}
}
