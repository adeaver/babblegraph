package usermetrics

import (
	"babblegraph/model/admin"
	"babblegraph/model/users"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "usermetrics",
	Routes: []router.Route{
		{
			Path: "get_user_aggregation_by_status_1",
			Handler: middleware.WithPermission(
				admin.PermissionViewUserMetrics,
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

func getUserAggregationByStatus(adminID admin.ID, r *router.Request) (interface{}, error) {
	resp := getUserAggregationByStatusResponse{}
	var aggregations []users.UserStatusAggregation
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		aggregations, err = users.GetUserStatusAggregation(tx)
		return err
	}); err != nil {
		return nil, err
	}
	for _, statusCount := range aggregations {
		switch statusCount.Status {
		case users.UserStatusVerified:
			resp.VerifiedUserCount = statusCount.Count
		case users.UserStatusUnverified:
			resp.UnverifiedUserCount = statusCount.Count
		case users.UserStatusUnsubscribed:
			resp.UnsubscribedUserCount = statusCount.Count
		case users.UserStatusBlocklistBounced,
			users.UserStatusBlocklistComplaint:
			resp.BlocklistedUserCount = statusCount.Count
		}
	}
	return resp, nil
}
