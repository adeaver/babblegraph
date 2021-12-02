package usermetrics

import (
	"babblegraph/admin/model/user"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
)

var Routes = router.RouteGroup{
	Prefix: "usermetrics",
	Routes: []router.Route{
		{
			Path: "get_user_aggregation_by_status_1",
			Handler: middleware.WithPermission(
				user.PermissionViewUserMetrics,
				getUserAggregationByStatus,
			),
		},
	},
}

type getUserAggregationByStatusResponse struct {
	VerifiedUserCount     int64 `json:"verified_user_count"`
	UnsubscribedUserCount int64 `json:"unsubscribed_user_count"`
	UnverifiedUserCount   int64 `json:"unverified_user_count"`
	BlocklistedUserCount  int64 `json:"blocklisted_user_count"`
}

func getUserAggregationByStatus(adminID user.AdminID, r *router.Request) (interface{}, error) {
	return getUserAggregationByStatusResponse{}, nil
}
