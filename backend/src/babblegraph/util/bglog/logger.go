package bglog

import (
	"log"

	"github.com/getsentry/sentry-go"
)

type Logger struct {
	contextKey    string
	sentryHub     *sentry.Hub
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	warningLogger *log.Logger
}
