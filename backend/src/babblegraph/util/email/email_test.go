package email

import (
	"babblegraph/util/testutils"
	"fmt"
	"testing"
)

func TestValidateEmailAddress(t *testing.T) {
	type testCase struct {
		input    string
		expected error
	}
	testCases := []testCase{
		{
			input:    "m@",
			expected: fmt.Errorf("Email Address is too long or too short"),
		}, {
			input:    "((((((@invalid.example.org",
			expected: fmt.Errorf("Invalid email address formatting"),
		}, {
			input:    "a-real-email-address@gmail.com",
			expected: nil,
		},
	}
	for idx, tc := range testCases {
		result := ValidateEmailAddress(tc.input)
		if err := testutils.CompareErrors(result, tc.expected); err != nil {
			t.Errorf("Error on test case %d: %s", idx+1, err.Error())
		}
	}
}
