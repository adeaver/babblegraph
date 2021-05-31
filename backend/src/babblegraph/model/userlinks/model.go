package userlinks

import (
	"babblegraph/model/email"
	"babblegraph/model/users"
	"fmt"
	"time"
)

type userLinkClickID string

type dbUserLinkClick struct {
	ID              userLinkClickID `db:"_id"`
	UserID          users.UserID    `db:"user_id"`
	Domain          string          `db:"domain"`
	URLIdentifier   string          `db:"url_identifier"`
	EmailRecordID   email.ID        `db:"email_record_id"`
	AccessMonth     AccessMonth     `db:"access_month"`
	FirstAccessedAt time.Time       `db:"first_accessed_at"`
}

type AccessMonth string

func getCurrentAccessMonth() string {
	nowUTC := time.Now().UTC()
	return fmt.Sprintf("%02d%d", nowUTC.Month(), nowUTC.Year())
}

type UserDomainCount struct {
	UserID users.UserID `db:"user_id"`
	Domain string       `db:"domain"`
	Count  int64        `db:"count"`
}

type paywallReportID string

type dbPaywallReport struct {
	ID            paywallReportID `db:"_id"`
	UserID        users.UserID    `db:"user_id"`
	Domain        string          `db:"domain"`
	URLIdentifier string          `db:"url_identifier"`
	EmailRecordID email.ID        `db:"email_record_id"`
	AccessMonth   AccessMonth     `db:"access_month"`
}
