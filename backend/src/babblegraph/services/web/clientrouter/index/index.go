package index

import (
	"babblegraph/model/utm"
	"babblegraph/services/web/clientrouter/initialdata"
	"babblegraph/services/web/clientrouter/middleware"
	"babblegraph/util/database"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"text/template"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func HandleServeIndexPage(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.LogRequestWithoutBody(r)
		handleUTMParameters(w, r)
		serveIndexTemplate(fmt.Sprintf("%s/index.html", staticFileDirName), w, r)
	}
}

func serveIndexTemplate(templateFileName string, w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(fmt.Sprintf("Error on index page: %s", err.Error()))
			sentry.CaptureException(err)
		}
		if x := recover(); x != nil {
			_, fn, line, _ := runtime.Caller(1)
			err := fmt.Errorf("Panic handling index: %s: %d: %v\n", fn, line, x)
			log.Println(fmt.Sprintf("Error on index page: %s", err.Error()))
			sentry.CaptureException(err)
		}
	}()
	w.Header().Add("Content-Type", "text/html")
	var tmpl *template.Template
	tmpl, err = template.New("index.html").ParseFiles(templateFileName)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	initialData, err := initialdata.GetInitialFrontendData(r)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = tmpl.Execute(w, *initialData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func handleUTMParameters(w http.ResponseWriter, r *http.Request) {
	utmParameters := utm.GetParametersForRequest(r)
	if utmParameters != nil {
		trackingID, err := utm.GetTrackingIDForRequest(r)
		switch {
		case err != nil:
			sentry.CaptureException(err)
		case trackingID != nil:
			localHub := sentry.CurrentHub().Clone()
			localHub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("utm-params", "go-func")
			})
			go func() {
				defer func() {
					if x := recover(); x != nil {
						_, fn, line, _ := runtime.Caller(1)
						err := fmt.Errorf("Panic handling utm parameters: %s: %d: %v\n", fn, line, x)
						localHub.CaptureException(err)
					}
				}()
				if err := database.WithTx(func(tx *sqlx.Tx) error {
					return utm.RegisterUTMPageHit(tx, *trackingID, *utmParameters)
				}); err != nil {
					localHub.CaptureException(err)
				}
			}()
			http.SetCookie(w, &http.Cookie{
				Name:     utm.UTMTrackingIDCookieName,
				Value:    trackingID.Str(),
				Secure:   true,
				HttpOnly: true,
				MaxAge:   int(utm.UTMTrackingIDMaxAge / time.Second),
				Expires:  time.Now().Add(utm.UTMTrackingIDMaxAge),
				Path:     "/",
			})
		default:
			log.Println("Unknown error making tracking ID for request, continuing...")
		}
	}
}
