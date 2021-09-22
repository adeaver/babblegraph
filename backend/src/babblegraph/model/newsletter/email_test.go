package newsletter

import (
	"babblegraph/model/email"
	"babblegraph/model/users"
)

type testEmailAccessor struct {
	emailRecords map[users.UserID][]email.ID
}

func (t *testEmailAccessor) InsertEmailRecord(id email.ID, userID users.UserID) error {
	emailRecordsForUser, _ := t.emailRecords[userID]
	t.emailRecords[userID] = append(emailRecordsForUser, id)
	return nil
}
