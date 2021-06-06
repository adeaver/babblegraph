package routes

import (
	"babblegraph/model/email"
	"babblegraph/model/users"
)

type ArticleLinkBodyDEPRECATED struct {
	UserID        users.UserID `json:"user_id"`
	EmailRecordID email.ID     `json:"email_record_id"`
	URL           string       `json:"url"`
}

type PaywallReportBodyDEPRECATED struct {
	UserID        users.UserID `json:"user_id"`
	EmailRecordID email.ID     `json:"email_record_id"`
	URL           string       `json:"url"`
}
