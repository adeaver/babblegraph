package tasks

import (
	"babblegraph/admin/model/user"
	"babblegraph/model/routes"
	"babblegraph/util/database"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func CreateAdminAndEmitToken(emailAddress string) error {
	return database.WithTx(func(tx *sqlx.Tx) error {
		if err := user.CreateAdminUser(tx, emailAddress); err != nil {
			return err
		}
		adminUser, err := user.LookupAdminUserByEmailAddress(tx, emailAddress)
		switch {
		case err != nil:
			return err
		case adminUser == nil:
			return fmt.Errorf("Not created")
		}
		token, err := routes.MakeAdminRegistrationToken(adminUser.AdminID)
		switch {
		case err != nil:
			return err
		case token == nil:
			return fmt.Errorf("Not created")
		}
		log.Println(*token)
		return nil
	})
}
