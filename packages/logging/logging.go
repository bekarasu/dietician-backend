package logging

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	fieldKeyMsg     = "message"
	fieldKeyTime    = "timestamp"
	timestampFormat = "2006-01-02 15:04:05"
)

type Config struct {
	Service ServiceConfig
}

func NewLogger(config Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	logger.SetFormatter(jsonFormatter())

	logger.AddHook(NewServiceHook(config.Service))
	logger.AddHook(NewTransactionHook())

	return logger
}

func jsonFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:  fieldKeyMsg,
			logrus.FieldKeyTime: fieldKeyTime,
		},
		TimestampFormat: timestampFormat,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// Return an empty string for the function name to remove it from the log output.
			// Keep the file name and line number intact.
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	}
}
