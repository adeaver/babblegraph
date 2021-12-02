package users

import (
	"babblegraph/util/math/decimal"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getUserStatusAggregationQuery = "SELECT status, COUNT(*) count FROM users GROUP BY status"

	getNetChangeAggregationQuery = "SELECT status, COUNT(*) count FROM users WHERE last_modified_at > current_date - interval '1 %s' GROUP BY status"
)

type UserStatusData struct {
	CurrentAggregation             []UserStatusAggregation
	VerifiedNetChangeOverLastWeek  decimal.Number
	VerifiedNetChangeOverLastMonth decimal.Number
}

func GetUserStatusData(tx *sqlx.Tx) (*UserStatusData, error) {
	userStatusAggregation, err := getUserStatusAggregation(tx)
	if err != nil {
		return nil, err
	}
	netChangeOverLastWeek, err := getVerifiedNetChangeOverPeriod(tx, "week")
	if err != nil {
		return nil, err
	}
	netChangeOverLastMonth, err := getVerifiedNetChangeOverPeriod(tx, "month")
	if err != nil {
		return nil, err
	}
	return &UserStatusData{
		CurrentAggregation:             userStatusAggregation,
		VerifiedNetChangeOverLastWeek:  *netChangeOverLastWeek,
		VerifiedNetChangeOverLastMonth: *netChangeOverLastMonth,
	}, nil
}

type UserStatusAggregation struct {
	Status UserStatus `db:"status"`
	Count  int64      `db:"count"`
}

func getUserStatusAggregation(tx *sqlx.Tx) ([]UserStatusAggregation, error) {
	var out []UserStatusAggregation
	if err := tx.Select(&out, getUserStatusAggregationQuery); err != nil {
		return nil, err
	}
	return out, nil
}

func getVerifiedNetChangeOverPeriod(tx *sqlx.Tx, period string) (*decimal.Number, error) {
	var agg []UserStatusAggregation
	if err := tx.Select(&agg, fmt.Sprintf(getNetChangeAggregationQuery, period)); err != nil {
		return nil, err
	}
	var netChange decimal.Number
	for _, status := range agg {
		value := decimal.FromInt64(status.Count)
		switch status.Status {
		case UserStatusVerified:
			netChange = netChange.Add(value)
		case UserStatusUnsubscribed,
			UserStatusBlocklistComplaint,
			UserStatusBlocklistBounced:
			netChange = netChange.Subtract(value)
		case UserStatusUnverified:
			// no-op
		}
	}
	return &netChange, nil
}
