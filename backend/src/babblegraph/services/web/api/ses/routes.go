package ses

import "babblegraph/services/web/router"

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "ses",
		Routes: []router.Route{
			{
				Path:    "handle_bounce_notification_1",
				Handler: handleBounceNotification,
			}, {
				Path:    "handle_complaint_notification_1",
				Handler: handleComplaintNotification,
			},
		},
	})
}

func handleBounceNotification(body []byte) (interface{}, error) {
	return nil, nil
}

func handleComplaintNotification(body []byte) (interface{}, error) {
	return nil, nil
}
