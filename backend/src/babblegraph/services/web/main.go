package main

import (
	"babblegraph/actions/email"
	"babblegraph/actions/verification"
	"babblegraph/model/routes"
	"babblegraph/model/utm"
	"babblegraph/services/web/api/language"
	"babblegraph/services/web/api/ses"
	"babblegraph/services/web/api/token"
	"babblegraph/services/web/api/user"
	utm_routes "babblegraph/services/web/api/utm"
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
	"github.com/jmoiron/sqlx"
)

func main() {
	log.Println("Starting babblegraph web server")
	r := mux.NewRouter()
	staticFileDirName := env.MustEnvironmentVariable("STATIC_DIR")

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

	// TODO: move these out of of here
	r.HandleFunc("/verify/{token}", func(w http.ResponseWriter, r *http.Request) {
		router.LogRequestWithoutBody(r)
		routeVars := mux.Vars(r)
		token, ok := routeVars["token"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userID, err := verification.VerifyUserByToken(token)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		subscriptionManagementLink, err := routes.MakeSubscriptionManagementRouteForUserID(*userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, *subscriptionManagementLink, http.StatusMovedPermanently)
	})
	// Warning: this next function is like one big hack.
	r.HandleFunc("/dist/{token}/logo.png", func(w http.ResponseWriter, r *http.Request) {
		router.LogRequestWithoutBody(r)
		// In order to collect information about whether an email was opened, we pass
		// the logo hero image with a URL of the above format. This is done because
		// 1) Some clients ban zero width images
		// 2) I imagine some clients ban url parameters
		routeVars := mux.Vars(r)
		go func() {
			// This is done in a go routine because we do not want it to
			// affect the speed with which requests are handled.
			token, ok := routeVars["token"]
			if !ok {
				return
			}
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return email.HandleDailyEmailOpenToken(tx, token)
			}); err != nil {
				log.Println(fmt.Sprintf("Got error handling token %s: %s", token, err.Error()))
			}
		}()
		http.ServeFile(w, r, fmt.Sprintf("%s/logo.png", staticFileDirName))
	})
	r.PathPrefix("/dist").Handler(http.StripPrefix("/dist", http.FileServer(http.Dir(staticFileDirName))))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.LogRequestWithoutBody(r)
		utmParameters := utm.GetParametersForRequest(r)
		if utmParameters != nil {
			trackingID, err := utm.GetTrackingIDForRequest(r)
			switch {
			case err != nil:
				log.Println(fmt.Sprintf("Error on getting tracking ID for request: %s", err.Error()))
			case trackingID != nil:
				go func() {
					if err := database.WithTx(func(tx *sqlx.Tx) error {
						return utm.RegisterUTMPageHit(tx, *trackingID, *utmParameters)
					}); err != nil {
						log.Println(fmt.Sprintf("Error registering UTM page hit: %s", err.Error()))
					}
				}()
				http.SetCookie(w, &http.Cookie{
					Name:  utm.UTMTrackingIDCookieName,
					Value: trackingID.Str(),
				})
			default:
				log.Println("Unknown error making tracking ID for request, continuing...")
			}
		}
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", staticFileDirName))
	})

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
	return nil
}
