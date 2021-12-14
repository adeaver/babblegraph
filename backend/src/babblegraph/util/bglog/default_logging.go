package bglog

import (
	"babblegraph/util/env"
	"fmt"
	"log"
)

// This package should eventually be deprecated in favor of context-based
// loggers, but remains here for backwards compatability

var (
	environmentName *env.Environment

	defaultDebugLogger   *log.Logger
	defaultInfoLogger    *log.Logger
	defaultWarningLogger *log.Logger
	defaultErrorLogger   *log.Logger
)

func InitLogger() {
	environmentName = env.MustEnvironmentName().Ptr()

	defaultDebugLogger = createLoggerForType(loggerTypeDebug)
	defaultInfoLogger = createLoggerForType(loggerTypeInfo)
	defaultWarningLogger = createLoggerForType(loggerTypeWarning)
	defaultErrorLogger = createLoggerForType(loggerTypeError)
}

func Debugf(format string, args ...interface{}) {
	switch *environmentName {
	case env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		defaultDebugLogger.Println(fmt.Sprintf(format, args...))
	case env.EnvironmentStage,
		env.EnvironmentProd:
		// no-op
	}
}

func Infof(format string, args ...interface{}) {
	defaultInfoLogger.Println(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	defaultWarningLogger.Println(fmt.Sprintf(format, args...))
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

func Errorf(format string, args ...interface{}) {
	defaultErrorLogger.Println(fmt.Sprintf(format, args...))
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

func Fatalf(format string, args ...interface{}) {
	log.Fatal(fmt.Sprintf(format, args...))
}
