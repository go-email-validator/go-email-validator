package log

import "github.com/sirupsen/logrus"

// Default logger logs to console by default, but user can replace the logger using the SetLibraryLogger()
var logger logrus.FieldLogger

func init() {
	logrus.StandardLogger().SetLevel(logrus.ErrorLevel)
	SetLogger(logrus.StandardLogger())
}

func SetLogger(l logrus.FieldLogger) {
	logger = l
}

func Logger() logrus.FieldLogger {
	return logger
}
