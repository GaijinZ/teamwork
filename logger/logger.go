package logger

import "github.com/sirupsen/logrus"

type Logger struct {
	Log *logrus.Logger
}

func NewLogger() *Logger {
	logger := logrus.New()
	return &Logger{Log: logger}
}

func (l *Logger) Infof(message string, args ...interface{}) {
	l.Log.Infof(message, args...)
}

func (l *Logger) Warningf(message string, args ...interface{}) {
	l.Log.Warnf(message, args...)
}

func (l *Logger) Errorf(message string, args ...interface{}) {
	l.Log.Errorf(message, args...)
}
