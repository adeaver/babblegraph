package bglog

import (
	"babblegraph/util/env"
	"fmt"
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

func NewLoggerForContext(tag, contextKey string) *Logger {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(tag, contextKey)
	})
	return &Logger{
		contextKey:    contextKey,
		sentryHub:     localHub,
		debugLogger:   createLoggerForType(loggerTypeDebug),
		infoLogger:    createLoggerForType(loggerTypeInfo),
		warningLogger: createLoggerForType(loggerTypeWarning),
		errorLogger:   createLoggerForType(loggerTypeError),
	}
}

func (l *Logger) logLineWithContext(format string, args ...interface{}) string {
	return fmt.Sprintf("%s | %s", l.contextKey, fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	switch env.MustEnvironmentName() {
	case env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		l.debugLogger.Println(l.logLineWithContext(format, args...))
	case env.EnvironmentStage,
		env.EnvironmentProd:
		// no-op
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.infoLogger.Println(l.logLineWithContext(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	logLineWithContext := l.logLineWithContext(format, args...)
	l.warningLogger.Println(logLineWithContext)
	switch *environmentName {
	case env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		// no-op
	case env.EnvironmentStage,
		env.EnvironmentProd:
		// TODO: add sentry here
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	logLineWithContext := l.logLineWithContext(format, args...)
	l.errorLogger.Println(logLineWithContext)
	switch *environmentName {
	case env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		// no-op
	case env.EnvironmentStage,
		env.EnvironmentProd:
		// TODO: add sentry here
	}
}
