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

func Info(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Info(msg)
	} else {
		Logger.Info(msg)
	}
}

func Error(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Error(msg)
	} else {
		Logger.Error(msg)
	}
}

func Debug(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Debug(msg)
	} else {
		Logger.Debug(msg)
	}
}

func Warn(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Warn(msg)
	} else {
		Logger.Warn(msg)
	}
}

func Fatal(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Fatal(msg)
	} else {
		Logger.Fatal(msg)
	}
}
