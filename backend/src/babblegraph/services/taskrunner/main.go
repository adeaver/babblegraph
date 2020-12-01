package main

import (
	"babblegraph/services/taskrunner/tasks"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"flag"
	"fmt"
	"log"
)

func main() {
	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	taskName := flag.String("task", "", "Name of task to run [daily-email]")
	if taskName == nil {
		log.Fatal("No task specified")
	}
	switch *taskName {
	case "daily-email":
		if err := tasks.SendDailyEmail(); err != nil {
			log.Fatal(err.Error())
		}
	default:
		log.Fatal(fmt.Sprintf("Invalid task specified %s", *taskName))
	}
}

func setupDatabases() error {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		return fmt.Errorf("Error setting up main-db: %s", err.Error())
	}
	if err := elastic.InitializeElasticsearchClientForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up elasticsearch: %s", err.Error())
	}
	return nil
}
