package advertising

import (
	"testing"
	"time"
)

func TestCalculateWeights(t *testing.T) {
	type testCase struct {
		input          []dbUserAdvertisement
		expectedOutput map[CampaignID]int
	}
	tcs := []testCase{
		{
			input: []dbUserAdvertisement{
				{
					CreatedAt:  time.Now().Add(-7 * 24 * time.Hour),
					CampaignID: CampaignID("1"),
				}, {
					CreatedAt:  time.Now().Add(-20 * 7 * 24 * time.Hour),
					CampaignID: CampaignID("2"),
				}, {
					CreatedAt:  time.Now().Add(-40 * 7 * 24 * time.Hour),
					CampaignID: CampaignID("2"),
				}, {
					CreatedAt:  time.Now().Add(-5 * 7 * 24 * time.Hour),
					CampaignID: CampaignID("3"),
				},
			},
			expectedOutput: map[CampaignID]int{
				CampaignID("1"): 0,
				CampaignID("2"): 10,
				CampaignID("3"): 1,
			},
		},
	}
	for idx, tc := range tcs {
		result := calculateWeightsForSeenAds(tc.input)
		for campaignID, weight := range result {
			expectedWeight, ok := tc.expectedOutput[campaignID]
			switch {
			case !ok:
				t.Errorf("Error on test case %d: result had campaign %s, but expected did not", idx+1, campaignID)
			case expectedWeight != weight:
				t.Errorf("Error on test case %d: got weight of %d for campaign id %s, but expected %d", idx+1, weight, campaignID, expectedWeight)
			}
			delete(tc.expectedOutput, campaignID)
		}
		if len(tc.expectedOutput) > 0 {
			for campaignID := range tc.expectedOutput {
				t.Errorf("Error on test case %d: expected output had campaign %s, but result did not", idx+1, campaignID)
			}
		}
	}
}

func TestGetWeight(t *testing.T) {
	type testCase struct {
		input map[CampaignID]int
	}
	tcs := []testCase{
		{
			input: map[CampaignID]int{
				CampaignID("1"): 5,
				CampaignID("2"): 0,
				CampaignID("3"): 2,
			},
		},
	}
	for idx, tc := range tcs {
		result := getWeightedCampaignIDs(tc.input)
		for _, campaignID := range result {
			tc.input[campaignID] = tc.input[campaignID] - 1
		}
		for campaignID, count := range tc.input {
			if count != 0 {
				t.Errorf("Error on test case %d: Campaign %s had a count of %d, but expected 0", idx+1, campaignID, count)
			}
		}
	}
}
