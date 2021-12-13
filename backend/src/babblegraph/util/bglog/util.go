package bglog

import (
	"fmt"
	"log"
	"strings"
)

type loggerType string

const (
	loggerTypeDebug   loggerType = "debug"
	loggerTypeInfo    loggerType = "info"
	loggerTypeWarning loggerType = "warning"
	loggerTypeError   loggerType = "error"
)

func (l loggerType) Str() string {
	return string(l)
}

func createLoggerForType(l loggerType) *log.Logger {
	return log.New(log.Writer(), fmt.Sprintf("%s: ", strings.ToUpper(l.Str())), log.Ldate|log.Ltime|log.Lshortfile)
}
