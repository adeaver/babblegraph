package ses

import (
	"babblegraph/model/users"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
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
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		for _, recipient := range req.Bounce.BouncedRecipients {
			didUpdate, err := users.AddUserToBlocklistByEmailAddress(tx, recipient.EmailAddress, users.UserStatusBlocklistBounced)
			if err != nil {
				log.Println(fmt.Sprintf("Error on adding %s to bounce list: %s", recipient.EmailAddress, err.Error()))
				continue
			}
			if didUpdate {
				log.Println(fmt.Sprintf("Successfully added %s to bounce list", recipient.EmailAddress))
			} else {
				log.Println(fmt.Sprintf("Did not add %s to bounce list", recipient.EmailAddress))
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return sesNotificationResponse{}, nil
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
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		for _, recipient := range req.Complaint.ComplainedRecipients {
			didUpdate, err := users.AddUserToBlocklistByEmailAddress(tx, recipient.EmailAddress, users.UserStatusBlocklistComplaint)
			if err != nil {
				log.Println(fmt.Sprintf("Error on adding %s to complaint list: %s", recipient.EmailAddress, err.Error()))
				continue
			}
			if didUpdate {
				log.Println(fmt.Sprintf("Successfully added %s to complaint list", recipient.EmailAddress))
			} else {
				log.Println(fmt.Sprintf("Did not add %s to complaint list", recipient.EmailAddress))
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return sesNotificationResponse{}, nil
}
