package useraccountsnotifications

import (
	"babblegraph/model/users"
	"babblegraph/util/ptr"
	"time"
)

type NotificationRequestID string

type dbNotificationRequest struct {
	ID             NotificationRequestID `db:"_id"`
	CreatedAt      time.Time             `db:"created_at"`
	LastModifiedAt time.Time             `db:"last_modified_at"`
	Type           NotificationType      `db:"notification_type"`
	UserID         users.UserID          `db:"user_id"`
	HoldUntil      time.Time             `db:"hold_until"`
	FulfilledAt    *time.Time            `db:"fulfilled_at"`
}

func (d dbNotificationRequest) ToNonDB() NotificationRequest {
	return NotificationRequest{
		ID:          d.ID,
		Type:        d.Type,
		UserID:      d.UserID,
		IsFulfilled: d.FulfilledAt != nil,
	}
}

type NotificationRequest struct {
	ID          NotificationRequestID
	Type        NotificationType
	UserID      users.UserID
	IsFulfilled bool
}

type NotificationType string

const (
	NotificationTypePaymentError                NotificationType = "payment_error"
	NotificationTypePremiumSubscriptionCanceled NotificationType = "premium_subscription_canceled"

	NotificationTypeNeedPaymentMethodWarning           NotificationType = "need_payment_method_warning"
	NotificationTypeTrialEndingSoon                    NotificationType = "trial_ending_soon"
	NotificationTypeNeedPaymentMethodWarningUrgent     NotificationType = "need_payment_method_warning_urgent"
	NotificationTypeNeedPaymentMethodWarningVeryUrgent NotificationType = "need_payment_method_warning_very_urgent"

	NotificationTypeAccountCreatedDEPRECATED            NotificationType = "account_created"
	NotificationTypeInitialPremiumInformationDEPRECATED NotificationType = "initial_premium_information"
)

type notificationRequestDebounceFulfillmentRecordID string

type dbNotificationRequestDebounceFulfillmentRecord struct {
	ID                    notificationRequestDebounceFulfillmentRecordID `db:"_id"`
	NotificationRequestID NotificationRequestID                          `db:"notification_request_id"`
	CreatedAt             time.Time                                      `db:"created_at"`
	LastModifiedAt        time.Time                                      `db:"created_at"`
}

// This defines the minimum time we need to wait between enqueuing messages of the same type
// If nil, only one message ever is allowed
var minimumElapsedTimeBetweenNotificationsByType = map[NotificationType]*time.Duration{
	NotificationTypeNeedPaymentMethodWarningVeryUrgent: nil,
	NotificationTypeNeedPaymentMethodWarning:           nil,
	NotificationTypeNeedPaymentMethodWarningUrgent:     nil,
	NotificationTypeTrialEndingSoon:                    nil,
	NotificationTypePaymentError:                       ptr.Duration(14 * 24 * time.Hour), // 2 weeks
	NotificationTypePremiumSubscriptionCanceled:        ptr.Duration(3 * 24 * time.Hour),  // 3 days
}
