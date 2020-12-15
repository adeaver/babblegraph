package ses

import (
	"babblegraph/services/web/router"
	"babblegraph/util/ses"
	"encoding/json"
	"fmt"
)

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

type sesNotificationResponse struct{}

func handleBounceNotification(body []byte) (interface{}, error) {
	var req ses.Notification
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	// TODO: persist the body here
	if req.Bounce == nil {
		return nil, fmt.Errorf("Bounce does not have a body")
	}
	return nil, nil
}

func handleComplaintNotification(body []byte) (interface{}, error) {
	var req ses.Notification
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	// TODO: persist the body here
	if req.Complaint == nil {
		return nil, fmt.Errorf("Complaint does not have a body")
	}
	return nil, nil
}
