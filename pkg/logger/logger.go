package logger

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init(level string) error {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})

	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	return nil
}

// Basic logging functions for simple use cases
func Info(msg string) {
	Logger.Info(msg)
}

func Error(msg string) {
	Logger.Error(msg)
}

func Fatal(msg string) {
	Logger.Fatal(msg)
}

// Unified error logging function for all error scenarios
func LogError(funcCtx, msg string, err error, fields ...logrus.Fields) {
	logFields := logrus.Fields{
		"context": funcCtx,
	}

	if err != nil {
		logFields["error"] = err.Error()
	}

	if len(fields) > 0 {
		for k, v := range fields[0] {
			logFields[k] = v
		}
	}

	Logger.WithFields(logFields).Error(msg)
}

// Success logging for positive outcomes
func LogSuccess(funcCtx, msg string, fields ...logrus.Fields) {
	logFields := logrus.Fields{"context": funcCtx}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			logFields[k] = v
		}
	}
	Logger.WithFields(logFields).Info(msg)
}
