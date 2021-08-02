package main

import (
	"babblegraph/services/web/api/language"
	"babblegraph/services/web/api/ses"
	"babblegraph/services/web/api/stripe"
	"babblegraph/services/web/api/token"
	"babblegraph/services/web/api/user"
	utm_routes "babblegraph/services/web/api/utm"
	"babblegraph/services/web/index"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting babblegraph web server")
	r := mux.NewRouter()

	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         env.MustEnvironmentVariable("SENTRY_DSN"),
		Environment: env.MustEnvironmentName().Str(),
	}); err != nil {
		log.Fatal(err.Error())
	}
	defer sentry.Flush(2 * time.Second)
	if err := registerAPI(r); err != nil {
		log.Fatal(err.Error())
	}
	if err := index.RegisterIndexRoutes(r, []index.IndexPage{}); err != nil {
		log.Fatal(err.Error())
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

func registerAPI(r *mux.Router) error {
	router.CreateNewAPIRouter(r)
	if err := user.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := ses.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := utm_routes.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := language.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := token.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := stripe.RegisterRouteGroups(); err != nil {
		return err
	}
	return nil
}
