package ses

import (
	"babblegraph/model/sesnotifications"
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
	switch {
	case req.SubscribeURL != nil:
		log.Println(fmt.Sprintf("Got subscription confirmation with URL: %s", *req.SubscribeURL))
		return nil, nil
	case req.Message != nil:
		var b ses.NotificationBody
		if err := json.Unmarshal([]byte(*req.Message), &b); err != nil {
			return nil, err
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := sesnotifications.InsertSESNotification(tx, req); err != nil {
				log.Println(fmt.Sprintf("Error persisting SES notification: %s. Continuing...", err.Error()))
			}
			for _, recipient := range b.Bounce.BouncedRecipients {
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
	default:
		return nil, fmt.Errorf("Bounce does not have a body")
	}
}

func handleComplaintNotification(body []byte) (interface{}, error) {
	var req ses.Notification
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	switch {
	case req.SubscribeURL != nil:
		log.Println(fmt.Sprintf("Got subscription confirmation with URL: %s", *req.SubscribeURL))
		return nil, nil
	case req.Message != nil:
		var b ses.NotificationBody
		if err := json.Unmarshal([]byte(*req.Message), &b); err != nil {
			return nil, err
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := sesnotifications.InsertSESNotification(tx, req); err != nil {
				log.Println(fmt.Sprintf("Error persisting SES notification: %s. Continuing...", err.Error()))
			}
			for _, recipient := range b.Complaint.ComplainedRecipients {
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
	default:
		return nil, fmt.Errorf("Complaint does not have a body")
	}
}
