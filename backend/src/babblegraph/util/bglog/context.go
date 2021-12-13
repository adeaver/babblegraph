package bglog

// If a function can potentially be used by web and async,
// then you can pass in a LogContext
type LogContext interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
