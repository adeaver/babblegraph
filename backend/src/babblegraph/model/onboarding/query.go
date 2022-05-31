package onboarding

import (
	"babblegraph/model/users"

	"github.com/jmoiron/sqlx"
)

func GetOnboardingStageForUser(tx *sqlx.Tx, userID users.UserID) (*OnboardingStage, error) {
	return OnboardingStageNotStarted.Ptr(), nil
}
