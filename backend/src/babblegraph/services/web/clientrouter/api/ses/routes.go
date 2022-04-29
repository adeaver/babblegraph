package ses

import (
	"babblegraph/model/billing"
	"babblegraph/model/email"
	"babblegraph/model/sesnotifications"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	email_util "babblegraph/util/email"
	"babblegraph/util/ses"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "ses",
	Routes: []router.Route{
		{
			Path:    "handle_bounce_notification_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(handleBounceNotification),
		}, {
			Path:    "handle_complaint_notification_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(handleComplaintNotification),
		},
	},
}

type sesNotificationResponse struct{}

func handleBounceNotification(r *router.Request) (interface{}, error) {
	var req ses.Notification
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	switch {
	case req.SubscribeURL != nil:
		r.Infof("Got subscription confirmation with URL: %s", *req.SubscribeURL)
		return sesNotificationResponse{}, nil
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
		for _, recipient := range b.Bounce.BouncedRecipients {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return handleAddUserToBlocklistByEmailAddress(r, tx, recipient.EmailAddress, users.UserStatusBlocklistBounced)
			}); err != nil {
				r.Errorf("Error adding user to blocklist: %s", err.Error())
				continue
			}
		}
		return sesNotificationResponse{}, nil
	default:
		return nil, fmt.Errorf("Bounce does not have a body")
	}
}

func handleComplaintNotification(r *router.Request) (interface{}, error) {
	var req ses.Notification
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	switch {
	case req.SubscribeURL != nil:
		r.Infof("Got subscription confirmation with URL: %s", *req.SubscribeURL)
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
		for _, recipient := range b.Complaint.ComplainedRecipients {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return handleAddUserToBlocklistByEmailAddress(r, tx, recipient.EmailAddress, users.UserStatusBlocklistComplaint)
			}); err != nil {
				r.Errorf("Error adding user to blocklist: %s", err.Error())
				continue
			}
		}
		return sesNotificationResponse{}, nil
	default:
		return nil, fmt.Errorf("Complaint does not have a body")
	}
}

func handleAddUserToBlocklistByEmailAddress(c ctx.LogContext, tx *sqlx.Tx, emailAddress string, blocklist users.UserStatus) error {
	formattedEmailAddress := email_util.FormatEmailAddress(emailAddress)
	user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
	switch {
	case err != nil:
		return err
	case user == nil:
		c.Infof("No user found for email address, skipping this user")
		return nil
	}
	shouldBlocklist, err := email.HandleBouncedEmail(tx, user.ID)
	switch {
	case err != nil:
		return err
	case !shouldBlocklist && blocklist != users.UserStatusBlocklistComplaint:
		c.Infof("User %s has been quarantined.", user.ID)
		return nil
	default:
		didUpdate, err := users.AddUserToBlocklist(tx, user.ID, blocklist)
		switch {
		case err != nil:
			return err
		case !didUpdate:
			c.Infof("Blocklisting user %s did not cause update", user.ID)
			return nil
		default:
			c.Infof("Added user %s to blocklist", user.ID)
			subscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
			switch {
			case err != nil:
				return fmt.Errorf("error cancelling subscription for user %s: %s", user.ID, err.Error())
			case subscription == nil:
				// no-op
			default:
				return billing.CancelPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
			}
		}
	}
	return nil
}
