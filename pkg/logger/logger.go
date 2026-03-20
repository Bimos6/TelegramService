package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(level string) {
	log = logrus.New()

	log.SetOutput(os.Stdout)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

func Get() *logrus.Logger {
	if log == nil {
		Init("info")
	}
	return log
}

func Debug(args ...interface{}) {
	Get().Debug(args...)
}

func Info(args ...interface{}) {
	Get().Info(args...)
}

func Warn(args ...interface{}) {
	Get().Warn(args...)
}

func Error(args ...interface{}) {
	Get().Error(args...)
}

func Fatal(args ...interface{}) {
	Get().Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	Get().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Get().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Get().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Get().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Get().Fatalf(format, args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return Get().WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Get().WithFields(fields)
}
