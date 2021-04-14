package usersurveys

import (
	"babblegraph/model/users"
	"time"
)

type SurveyID string

func (s SurveyID) Ptr() *SurveyID {
	return &s
}

type dbSurvey struct {
	UserID        users.UserID `db:"user_id"`
	SurveyID      SurveyID     `db:"_id"`
	SurveyType    SurveyType   `db:"survey_type"`
	SentAt        time.Time    `db:"sent_at"`
	FirstOpenedAt *time.Time   `db:"first_opened_at"`
}

type SurveyType string

const (
	SurveyTypeLowOpen1     SurveyType = "LowOpen1"
	SurveyTypeHighOpen1    SurveyType = "HighOpen1"
	SurveyTypeUnsubscribe1 SurveyType = "Unsubscribe1"
)

func (t SurveyType) Ptr() *SurveyType {
	return &t
}

type surveyResponseID string

type dbSurveyResponse struct {
	ID          surveyResponseID `db:"_id"`
	QuestionID  QuestionID       `db:"question_id"`
	SurveyID    SurveyID         `db:"survey_id"`
	Answer      string           `db:"answer"`
	SubmittedAt time.Time        `db:"submitted_at"`
}

type QuestionID string
