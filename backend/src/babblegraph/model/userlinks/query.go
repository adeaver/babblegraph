package userlinks

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

const (
	registerUserLinkClickQuery               = "INSERT INTO user_link_clicks (user_id, domain, source_id, url_identifier, email_record_id, access_month) VALUES ($1, $2, $3, $4, $5, $6)"
	getDomainCountsByCurrentAccessMonthQuery = "SELECT user_id, domain, COUNT(DISTINCT url_identifier) count FROM user_link_clicks WHERE user_id = $1 AND access_month = $2 GROUP BY user_id, domain"

	reportPaywallQuery = "INSERT INTO paywall_reports (user_id, domain, url_identifier, email_record_id, access_month) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id, url_identifier, access_month) DO NOTHING"
)

func RegisterUserLinkClick(tx *sqlx.Tx, userID users.UserID, u urlparser.ParsedURL, emailRecordID email.ID) error {
	currentAccessMonth := getCurrentAccessMonth()
	sourceID, err := content.GetSourceIDForParsedURL(tx, u)
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

// TODO(content-migration): remove this
func BackfillUserLinkClicks(c ctx.LogContext, tx *sqlx.Tx) error {
	rows, err := tx.Queryx("SELECT * FROM user_link_clicks WHERE source_id IS NULL")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var l dbUserLinkClick
		if err := rows.StructScan(&l); err != nil {
			return err
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			sourceID, err := content.LookupSourceIDForParsedURL(tx, urlparser.ParsedURL{
				Domain: l.Domain,
			})
			switch {
			case err != nil:
				return err
			case sourceID == nil:
				return nil
			default:
				_, err = tx.Exec("UPDATE user_link_clicks SET source_id = $1 WHERE _id = $2", *sourceID, l.ID)
				return err
			}
		}); err != nil {
			c.Infof("Error processing link click %s: %s", l.ID, err.Error())
		}
	}
	return nil
}
