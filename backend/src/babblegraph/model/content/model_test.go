package content

import "testing"

func TestTopicMappingIDOriginID(t *testing.T) {
	type testCase struct {
		SourceSeedTopicMappingID *SourceSeedTopicMappingID
		SourceTopicMappingID     *SourceTopicMappingID
	}
	tcs := []testCase{
		{
			SourceSeedTopicMappingID: SourceSeedTopicMappingID("12345-2345").Ptr(),
		}, {
			SourceTopicMappingID: SourceTopicMappingID("12345-2345").Ptr(),
		},
	}
	for idx, tc := range tcs {
		result := MustMakeTopicMappingID(MakeTopicMappingIDInput{
			SourceSeedTopicMappingID: tc.SourceSeedTopicMappingID,
			SourceTopicMappingID:     tc.SourceTopicMappingID,
		})
		sourceSeedTopicMappingID, sourceTopicMappingID, err := result.GetOriginID()
		switch {
		case err != nil:
			t.Errorf("Got error on test %d: %s", idx+1, err.Error())
		case tc.SourceSeedTopicMappingID != nil:
			switch {
			case sourceSeedTopicMappingID == nil:
				t.Errorf("Error on test %d: expected result to be source seed topic mapping id %s, but got null", idx+1, *tc.SourceSeedTopicMappingID)
			case *sourceSeedTopicMappingID != *tc.SourceSeedTopicMappingID:
				t.Errorf("Error on test %d: expected result to be source seed topic mapping id %s, but got %s", idx+1, *tc.SourceSeedTopicMappingID, *sourceSeedTopicMappingID)
			}
		case tc.SourceTopicMappingID != nil:
			switch {
			case sourceTopicMappingID == nil:
				t.Errorf("Error on test %d: expected result to be source topic mapping id %s, but got null", idx+1, *tc.SourceTopicMappingID)
			case *sourceTopicMappingID != *tc.SourceTopicMappingID:
				t.Errorf("Error on test %d: expected result to be source topic mapping id %s, but got %s", idx+1, *tc.SourceTopicMappingID, *sourceSeedTopicMappingID)
			}
		}
	}
}
