package main

import (
	"babblegraph/services/web/adminrouter"
	"babblegraph/services/web/clientrouter"
	"babblegraph/util/bglog"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/wordsmith"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

func main() {
	bglog.InitLogger()
	bglog.Infof("Starting babblegraph web server")
	r := mux.NewRouter()
	if err := setupDatabases(); err != nil {
		bglog.Fatalf("Error initializing databases: %s", err.Error())
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         env.MustEnvironmentVariable("SENTRY_DSN"),
		Environment: env.MustEnvironmentName().Str(),
	}); err != nil {
		bglog.Fatalf("Error initializing sentry: %s", err.Error())
	}
	defer sentry.Flush(2 * time.Second)
	if err := adminrouter.RegisterAdminRouter(r); err != nil {
		bglog.Fatalf("Error adding admin router: %s", err.Error())
	}
	if err := clientrouter.RegisterClientRouter(r); err != nil {
		bglog.Fatalf("Error adding client router router: %s", err.Error())
	}

	http.ListenAndServe(":8080", r)
}

func setupDatabases() error {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		return fmt.Errorf("Error setting up main-db: %s", err.Error())
	}
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up wordsmith: %s", err.Error())
	}
	return nil
}
