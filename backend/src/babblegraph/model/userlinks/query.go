package userlinks

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

const (
	registerUserLinkClickQuery               = "INSERT INTO user_link_clicks (user_id, domain, source_id, url_identifier, email_record_id, access_month) VALUES ($1, $2, $3, $4, $5, $6)"
	getDomainCountsByCurrentAccessMonthQuery = "SELECT user_id, source_id, COUNT(DISTINCT url_identifier) count FROM user_link_clicks WHERE user_id = $1 AND access_month = $2 AND source_id IS NOT NULL GROUP BY user_id, source_id"

	reportPaywallQuery = "INSERT INTO paywall_reports (user_id, domain, url_identifier, email_record_id, access_month) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id, url_identifier, access_month) DO NOTHING"
)

func RegisterUserLinkClick(tx *sqlx.Tx, userID users.UserID, u urlparser.ParsedURL, emailRecordID email.ID) error {
	currentAccessMonth := getCurrentAccessMonth()
	parsedDomain := urlparser.MustParseURL(u.Domain)
	sourceID, err := content.GetSourceIDForParsedURL(tx, parsedDomain)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(registerUserLinkClickQuery, userID, u.Domain, sourceID, u.URLIdentifier, emailRecordID, currentAccessMonth); err != nil {
		return err
	}
	return nil
}

func GetDomainCountsByCurrentAccessMonthForUser(tx *sqlx.Tx, userID users.UserID) ([]UserDomainCount, error) {
	currentAccessMonth := getCurrentAccessMonth()
	var domainCounts []UserDomainCount
	if err := tx.Select(&domainCounts, getDomainCountsByCurrentAccessMonthQuery, userID, currentAccessMonth); err != nil {
		return nil, err
	}
	return domainCounts, nil
}

func ReportPaywall(tx *sqlx.Tx, userID users.UserID, u urlparser.ParsedURL, emailRecordID email.ID) error {
	currentAccessMonth := getCurrentAccessMonth()
	if _, err := tx.Exec(reportPaywallQuery, userID, u.Domain, u.URLIdentifier, emailRecordID, currentAccessMonth); err != nil {
		return err
	}
	return nil
}
