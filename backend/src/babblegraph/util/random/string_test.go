package random

import (
	"babblegraph/util/testutils"
	"fmt"
	"testing"
)

func TestRandomStringCreation(t *testing.T) {
	type testCase struct {
		inputLength int
		expectedErr error
	}
	testCases := []testCase{
		{
			inputLength: 5,
			expectedErr: nil,
		}, {
			inputLength: 0,
			expectedErr: fmt.Errorf("Must have at least length of 1"),
		}, {
			inputLength: -1,
			expectedErr: fmt.Errorf("Must have at least length of 1"),
		}, {
			inputLength: 10,
			expectedErr: nil,
		},
	}
	for idx, tc := range testCases {
		s, err := MakeRandomString(tc.inputLength)
		if err := testutils.CompareErrors(tc.expectedErr, err); err != nil {
			t.Errorf("Error on test case %d: %s", idx+1, err.Error())
			continue
		}
		if len(*s) != tc.inputLength {
			t.Errorf("Error on test case %d: expected string of length %d, but got string of length %d", idx+1, tc.inputLength, len(*s))
		}
	}
}
