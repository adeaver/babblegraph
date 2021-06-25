package useraccounts

import (
	"babblegraph/model/users"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	createUserPasswordQuery = "INSERT INTO user_account_passwords (password_hash, user_id, salt) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET password_hash = $1, salt = $3, created_at = timezone('utc', now())"
	getPasswordForUserQuery = "SELECT * FROM user_account_passwords WHERE user_id = $1"

	hasSubscriptionForUserQuery      = "SELECT * FROM user_account_subscription_levels WHERE user_id = $1"
	getSubscriptionLevelForUserQuery = "SELECT * FROM user_account_subscription_levels WHERE user_id = $1 AND is_active = TRUE"
	addSubscriptionLevelForUserQuery = "INSERT INTO user_account_subscription_levels (user_id, subscription_level, expires_at) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET is_active = TRUE, subscription_level = $2, expires_at = $3"
	expireSubscriptionForUserQuery   = "UPDATE user_account_subscription_levels SET is_active = FALSE WHERE user_id = $1"

	forgotPasswordExpirationTime   = 15 * 60 * time.Second
	maxDailyForgotPasswordRequests = 5

	getAllUnfulfilledForgotPasswordAttemptsQuery       = "SELECT * FROM user_forgot_password_attempts WHERE is_archived IS FALSE AND fulfilled_at IS NULL"
	fulfillForgotPasswordAttemptByIDQuery              = "UPDATE user_forgot_password_attempts SET fulfilled_at = (now() at time zone 'utc') WHERE _id = $1"
	getNotArchivedForgotPasswordAttemptsForUserIDQuery = "SELECT * FROM user_forgot_password_attempts WHERE is_archived IS FALSE AND user_id = $1"
	getNotArchivedForgotPasswordAttemptsForIDQuery     = "SELECT * FROM user_forgot_password_attempts WHERE is_archived IS FALSE AND _id = $1"
	addForgotPasswordAttemptQuery                      = "INSERT INTO user_forgot_password_attempts (user_id) VALUES ($1) ON CONFLICT user_id DO NOTHING"
	archiveAllFulfilledForgotPasswordAttemptsQuery     = "UPDATE user_forgot_password_attempts SET is_archived = TRUE WHERE is_archived = FALSE AND fulfilled_at IS NOT NULL AND fulfilled_at < NOW() - interval '20 minutes'"
)

func CreateUserPasswordForUser(tx *sqlx.Tx, userID users.UserID, password string) error {
	if !ValidatePasswordMeetsRequirements(password) {
		return fmt.Errorf("password does not meet requirements")
	}
	passwordSalt, err := generatePasswordSalt()
	if err != nil {
		return err
	}
	passwordHash, err := generatePasswordHash(password, *passwordSalt)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(createUserPasswordQuery, *passwordHash, userID, *passwordSalt); err != nil {
		return err
	}
	return nil
}

func DoesUserAlreadyHaveAccount(tx *sqlx.Tx, userID users.UserID) (bool, error) {
	var matches []dbUserPassword
	if err := tx.Select(&matches, getPasswordForUserQuery, userID); err != nil {
		return false, err
	}
	return len(matches) > 0, nil
}

func VerifyPasswordForUser(tx *sqlx.Tx, userID users.UserID, password string) error {
	var matches []dbUserPassword
	if err := tx.Select(&matches, getPasswordForUserQuery, userID); err != nil {
		return err
	}
	if len(matches) != 1 {
		return fmt.Errorf("Incorrect number of passwords matching found")
	}
	return comparePasswords(matches[0].PasswordHash, password, matches[0].Salt)
}

func DoesUserHaveSubscription(tx *sqlx.Tx, userID users.UserID) (_didHaveAccount, _alreadyHadActiveAccount bool, _err error) {
	var matches []dbUserSubscription
	if err := tx.Select(&matches, hasSubscriptionForUserQuery, userID); err != nil {
		return false, false, err
	}
	if len(matches) < 1 {
		return false, false, nil
	}
	return true, matches[0].IsActive, nil
}

func LookupSubscriptionLevelForUser(tx *sqlx.Tx, userID users.UserID) (*SubscriptionLevel, error) {
	var matches []dbUserSubscription
	if err := tx.Select(&matches, getSubscriptionLevelForUserQuery, userID); err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, nil
	}
	return matches[0].SubscriptionLevel.Ptr(), nil
}

func getDefaultExpirationTime() time.Time {
	// This is basically totally meaningless right now
	now := time.Now()
	return time.Date(now.Year()+1, now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location()).UTC()
}

func AddSubscriptionLevelForUser(tx *sqlx.Tx, userID users.UserID, level SubscriptionLevel) error {
	if _, err := tx.Exec(addSubscriptionLevelForUserQuery, userID, level, getDefaultExpirationTime()); err != nil {
		return err
	}
	return nil
}

func ExpireSubscriptionForUser(tx *sqlx.Tx, userID users.UserID) error {
	if _, err := tx.Exec(expireSubscriptionForUserQuery, userID); err != nil {
		return err
	}
	return nil
}

func GetAllUnfulfilledForgotPasswordAttempts(tx *sqlx.Tx) ([]ForgotPasswordAttempt, error) {
	var matches []dbUserForgotPasswordAttempt
	if err := tx.Select(&matches, getAllUnfulfilledForgotPasswordAttemptsQuery); err != nil {
		return nil, err
	}
	var out []ForgotPasswordAttempt
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func FulfillForgotPasswordAttempt(tx *sqlx.Tx, id ForgotPasswordAttemptID) error {
	if _, err := tx.Exec(fulfillForgotPasswordAttemptByIDQuery, id); err != nil {
		return err
	}
	return nil
}

func GetUnexpiredForgotPasswordAttemptByID(tx *sqlx.Tx, id ForgotPasswordAttemptID) (_res *ForgotPasswordAttempt, _isExpired bool, _err error) {
	// A forgot password attempt is considered unexpired if:
	// 1) It is not archived
	// 2) It was fulfilled less than the expiration time ago
	var matches []dbUserForgotPasswordAttempt
	if err := tx.Select(&matches, getNotArchivedForgotPasswordAttemptsForIDQuery); err != nil {
		return nil, false, err
	}
	switch {
	case len(matches) == 0:
		// This is not an error because it can be easily user generated
		// is expected if a user clicks a very old link and we want to send a sentry
		// if there's a real error
		log.Println(fmt.Sprintf("No unarchived password attempt found for ID: %s", id))
		return nil, false, nil
	case len(matches) > 1:
		return nil, false, fmt.Errorf("Expected 1 unarchived password attempt for ID %s, but got %d", id, len(matches))
	case matches[0].FulfilledAt == nil:
		return nil, false, fmt.Errorf("Forgot password attempt with ID %s is unfulfilled", id)
	default:
		expirationTime := matches[0].FulfilledAt.Add(forgotPasswordExpirationTime)
		if time.Now().After(expirationTime) {
			return nil, true, nil
		}
		out := matches[0].ToNonDB()
		return &out, false, nil
	}
}

func AddForgotPasswordAttemptForUserID(tx *sqlx.Tx, userID users.UserID) (_hasTooManyAttempts bool, _err error) {
	var matches []dbUserForgotPasswordAttempt
	if err := tx.Select(&matches, getNotArchivedForgotPasswordAttemptsForUserIDQuery); err != nil {
		return false, err
	}
	if len(matches) >= maxDailyForgotPasswordRequests {
		return true, nil
	}
	if _, err := tx.Exec(addForgotPasswordAttemptQuery, userID); err != nil {
		return false, err
	}
	return false, nil
}

func ArchiveAllForgotPasswordAttemptsOlderThan20Minutes(tx *sqlx.Tx) error {
	if _, err := tx.Exec(archiveAllFulfilledForgotPasswordAttemptsQuery); err != nil {
		return err
	}
	return nil
}
