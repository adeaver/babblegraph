package main

import (
	"log"

	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"babblegraph/worker/queuedefs"
)

func main() {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		log.Fatal(err.Error())
	}
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		log.Fatal(err.Error())
	}
	errs := make(chan error, 1)
	if err := queuedefs.RegisterQueues(errs); err != nil {
		log.Fatal(err.Error())
	}
	<-errs
}
