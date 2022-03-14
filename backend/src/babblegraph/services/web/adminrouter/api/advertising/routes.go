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
		},
	},
}
