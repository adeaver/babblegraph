package useraccounts

import (
	"babblegraph/model/users"
	"fmt"
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