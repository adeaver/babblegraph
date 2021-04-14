package usersurveys

import (
	"babblegraph/model/users"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	createSurveyForUserAndTypeQuery = "INSERT INTO user_surveys (_id, user_id, survey_type) VALUES ($1, $2, $3)"
	getSurveyForSurveyIDQuery       = "SELECT * FROM user_surveys WHERE _id = $1"
	updateSurveyFirstOpenedAtTime   = "UPDATE user_surveys SET first_opened_at = timezone('utc', now()) WHERE _id = $1 AND first_opened_at IS NULL"
	submitResponseForSurvey         = "INSERT INTO user_survey_responses (question_id, survey_id, answer) VALUES ($1, $2, $3)"
)

func CreateSurveyForUserAndType(tx *sqlx.Tx, userID users.UserID, surveyType SurveyType) (*SurveyID, error) {
	surveyID := SurveyID(uuid.New().String())
	if _, err := tx.Exec(createSurveyForUserAndTypeQuery, surveyID, userID, surveyType); err != nil {
		return nil, err
	}
	return surveyID.Ptr(), nil
}

func GetSurveyTypeForID(tx *sqlx.Tx, surveyID SurveyID) (*SurveyType, error) {
	var matches []dbSurvey
	if err := tx.Select(&matches, getSurveyForSurveyIDQuery, surveyID); err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, fmt.Errorf("Expected one survey match, got %d", len(matches))
	}
	if _, err := tx.Exec(updateSurveyFirstOpenedAtTime, surveyID); err != nil {
		return nil, err
	}
	return matches[0].SurveyType.Ptr(), nil
}

func SubmitResponseForSurvey(tx *sqlx.Tx, surveyID SurveyID, questionID QuestionID, answer string) error {
	if _, err := tx.Exec(submitResponseForSurvey, questionID, surveyID, answer); err != nil {
		return err
	}
	return nil
}
