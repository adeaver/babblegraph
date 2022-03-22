package advertising

import "fmt"

func getExperimentNameForCampaignID(id CampaignID) string {
	return fmt.Sprintf("advertising_campaign_%s", id)
}
