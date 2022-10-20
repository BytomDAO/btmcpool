package logger

import (
	"github.com/sirupsen/logrus"
)

type Level uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Convert the slice to logrus.Fields
func sliceToFields(params []interface{}) logrus.Fields {
	if len(params)%2 != 0 {
		logrus.Error("Log parameter length is wrong!")
		return nil
	}
	fields := logrus.Fields{}
	for i := 0; i < len(params); i += 2 {
		val := params[i]
		key, found := val.(string)
		if found {
			fields[key] = params[i+1]
		} else {
			return nil
		}
	}
	return fields
}
