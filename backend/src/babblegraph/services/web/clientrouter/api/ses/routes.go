package ses

import (
	"babblegraph/model/billing"
	"babblegraph/model/sesnotifications"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/api"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"encoding/json"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return api.RegisterRouteGroup(api.RouteGroup{
		Prefix: "ses",
		Routes: []api.Route{
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
			return sesnotifications.InsertSESNotification(tx, req)
		}); err != nil {
			return nil, err
		}
		c := ctx.GetDefaultLogContext()
		for _, recipient := range b.Bounce.BouncedRecipients {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				didUpdate, err := users.AddUserToBlocklistByEmailAddress(tx, recipient.EmailAddress, users.UserStatusBlocklistBounced)
				if err != nil {
					return fmt.Errorf("Error on adding %s to bounce list: %s", recipient.EmailAddress, err.Error())
				}
				if didUpdate {
					log.Println(fmt.Sprintf("Successfully added %s to bounce list", recipient.EmailAddress))
					user, err := users.LookupUserByEmailAddress(tx, recipient.EmailAddress)
					switch {
					case err != nil:
						return fmt.Errorf("Error finding user %s: %s", recipient.EmailAddress, err.Error())
					case user == nil:
						log.Println(fmt.Sprintf("No user found %s", recipient.EmailAddress))
					default:
						subscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
						switch {
						case err != nil:
							return fmt.Errorf("error cancelling subscription for user %s: %s", recipient.EmailAddress, err.Error())
						case subscription == nil:
							// no-op
						default:
							return billing.CancelPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
						}
					}
				} else {
					log.Println(fmt.Sprintf("Did not add %s to bounce list", recipient.EmailAddress))
				}
				return nil
			}); err != nil {
				sentry.CaptureException(err)
				continue
			}
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
			return sesnotifications.InsertSESNotification(tx, req)
		}); err != nil {
			sentry.CaptureException(fmt.Errorf("Error persisting SES notification: %s. Continuing...", err.Error()))
		}
		c := ctx.GetDefaultLogContext()
		for _, recipient := range b.Complaint.ComplainedRecipients {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				didUpdate, err := users.AddUserToBlocklistByEmailAddress(tx, recipient.EmailAddress, users.UserStatusBlocklistComplaint)
				if err != nil {
					return fmt.Errorf("Error on adding %s to complaint list: %s", recipient.EmailAddress, err.Error())
				}
				if didUpdate {
					log.Println(fmt.Sprintf("Successfully added %s to complaint list", recipient.EmailAddress))
					user, err := users.LookupUserByEmailAddress(tx, recipient.EmailAddress)
					switch {
					case err != nil:
						return fmt.Errorf("Error finding user %s: %s", recipient.EmailAddress, err.Error())
					case user == nil:
						log.Println(fmt.Sprintf("No user found %s", recipient.EmailAddress))
					default:
						subscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
						switch {
						case err != nil:
							return fmt.Errorf("error cancelling subscription for user %s: %s", recipient.EmailAddress, err.Error())
						case subscription == nil:
							// no-op
						default:
							return billing.CancelPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
						}
					}
				} else {
					log.Println(fmt.Sprintf("Did not add %s to complaint list", recipient.EmailAddress))
				}
				return nil
			}); err != nil {
				return nil, err
			}
		}
		return sesNotificationResponse{}, nil
	default:
		return nil, fmt.Errorf("Complaint does not have a body")
	}
}
