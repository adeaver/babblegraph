package advertising

import (
	"babblegraph/model/admin"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
)

var Routes = router.RouteGroup{
	Prefix: "advertising",
	Routes: []router.Route{
		{
			Path: "get_all_vendors_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingVendors,
				getAllVendors,
			),
		}, {
			Path: "insert_vendor_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingVendors,
				insertVendor,
			),
		}, {
			Path: "update_vendor_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingVendors,
				editVendor,
			),
		}, {
			Path: "get_all_sources_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingSources,
				getAllSources,
			),
		}, {
			Path: "insert_source_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingSources,
				insertSource,
			),
		}, {
			Path: "update_source_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingSources,
				editSource,
			),
		}, {
			Path: "get_all_campaigns_1",
			Handler: middleware.WithPermission(
				admin.PermissionViewAdvertisingCampaigns,
				getAllCampaigns,
			),
		}, {
			Path: "get_campaign_1",
			Handler: middleware.WithPermission(
				admin.PermissionViewAdvertisingCampaigns,
				getCampaign,
			),
		}, {
			Path: "insert_campaign_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingCampaigns,
				insertCampaign,
			),
		}, {
			Path: "update_campaign_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingCampaigns,
				updateCampaign,
			),
		}, {
			Path: "get_campaign_topic_mappings_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingCampaigns,
				getCampaignTopicMappings,
			),
		}, {
			Path: "update_campaign_topic_mappings_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingCampaigns,
				updateCampaignTopicMappings,
			),
		}, {
			Path: "get_all_advertisements_for_campaign_1",
			Handler: middleware.WithPermission(
				admin.PermissionViewAdvertisingAdvertisements,
				getAllAdvertisementsForCampaign,
			),
		}, {
			Path: "insert_advertisement_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingAdvertisements,
				insertAdvertisement,
			),
		}, {
			Path: "update_advertisement_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingAdvertisements,
				updateAdvertisement,
			),
		},
	},
}
