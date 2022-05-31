package onboarding

import (
	"babblegraph/model/users"
	"time"
)

type OnboardingStage string

const (
	OnboardingStageNotStarted        OnboardingStage = "not-started"
	OnboardingStageSettings          OnboardingStage = "settings"
	OnboardingStageInterestSelection OnboardingStage = "interest-selection"
	OnboardingStageVocabulary        OnboardingStage = "vocabulary"
	OnboardingStageFinished          OnboardingStage = "finished"
)

func (o OnboardingStage) Ptr() *OnboardingStage {
	return &o
}

type OnboardingID string

type dbOnboarding struct {
	CreatedAt      time.Time       `db:"created_at"`
	LastModifiedAt time.Time       `db:"last_modified_at"`
	ID             OnboardingID    `db:"_id"`
	UserID         users.UserID    `db:"user_id"`
	Status         OnboardingStage `db:"status"`
}
