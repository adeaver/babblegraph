package main

import (
	"babblegraph/worker/queuedefs"
	"babblegraph/worker/wordsmith"
	"log"

	"github.com/adeaver/babblegraph/lib/database"
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
