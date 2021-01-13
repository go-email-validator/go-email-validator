package log

import (
	"go.uber.org/zap"
)

// Default logger logs to console by default
var logger *zap.Logger

func init() {
	lBuilder := zap.NewProductionConfig()
	lBuilder.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	l, err := lBuilder.Build()
	if err != nil {
		panic(err)
	}
	SetLogger(l)
}

// SetLogger sets logger
func SetLogger(l *zap.Logger) {
	logger = l
}

// Logger returns logger
func Logger() *zap.Logger {
	return logger
}
