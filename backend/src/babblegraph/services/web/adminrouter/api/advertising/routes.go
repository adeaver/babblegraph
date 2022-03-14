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
			Path: "insert_campaign_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingCampaigns,
				insertCampaign,
			),
		}, {
			Path: "update_campaign_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditAdvertisingCampaigns,
				insertCampaign,
			),
		},
	},
}
