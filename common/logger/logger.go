package logger

import (
	officiallog "log"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	contextLogger *logrus.Logger
	entry         *logrus.Entry
}

func newLogger(fields logrus.Fields, logLevel Level) *Logger {
	// Create new logger
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000000",
	})

	// set log level
	log.SetLevel(logrus.Level(logLevel))

	return &Logger{
		contextLogger: log,
		entry:         log.WithFields(fields),
	}
}

// The global Logger
var logger *Logger

// Initialize the Logger
func Init(logLevel Level) error {
	return InitWithFields(logLevel, map[string]interface{}{})
}

// Initialize the Logger with fields
func InitWithFields(logLevel Level, fields map[string]interface{}) error {
	logger = newLogger(logrus.Fields{}, logLevel)
	WithFields(fields)
	return nil
}

func WithFields(fields map[string]interface{}) {
	logger.entry = logger.entry.WithFields(fields)
}

func Info(msg string, params ...interface{}) {
	fields := sliceToFields(params)
	logger.entry.WithFields(fields).Info(msg)
}

func Debug(msg string, params ...interface{}) {
	fields := sliceToFields(params)
	logger.entry.WithFields(fields).Debug(msg)
}

func Warn(msg string, params ...interface{}) {
	fields := sliceToFields(params)
	logger.entry.WithFields(fields).Warn(msg)
}

func Error(msg string, params ...interface{}) {
	fields := sliceToFields(params)
	logger.entry.WithFields(fields).Error(msg)
}

func Fatal(msg string, params ...interface{}) {
	fields := sliceToFields(params)
	logger.entry.WithFields(fields).Fatal(msg)
}

func Println(v ...interface{}) {
	officiallog.Println(v...)
}
