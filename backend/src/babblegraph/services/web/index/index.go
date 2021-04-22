package index

import (
	"babblegraph/model/utm"
	"babblegraph/util/database"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func ServeIndexPage(indexFilenameWithPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		utmParameters := utm.GetParametersForRequest(r)
		if utmParameters != nil {
			handleUTMParameters(*utmParameters, w, r)
		}
		serveIndexTemplate(indexFilenameWithPath, struct{}{}, w, r)
	}
}

func serveIndexTemplate(templateFileName string, templateData interface{}, w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			sentry.CaptureException(err)
		}
	}()
	fileNameParts := strings.Split(templateFileName, "/")
	var tmpl *template.Template
	tmpl, err = template.New(fileNameParts[len(fileNameParts)-1]).ParseFiles(templateFileName)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	if err := tmpl.Execute(w, templateData); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func handleUTMParameters(utmParameters utm.Parameters, w http.ResponseWriter, r *http.Request) {
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
				return utm.RegisterUTMPageHit(tx, *trackingID, utmParameters)
			}); err != nil {
				localHub.CaptureException(err)
			}
		}()
		http.SetCookie(w, &http.Cookie{
			Name:  utm.UTMTrackingIDCookieName,
			Value: trackingID.Str(),
		})
	default:
		err := fmt.Errorf("Unknown error making tracking ID for request, continuing...")
		sentry.CaptureException(err)
	}
}
