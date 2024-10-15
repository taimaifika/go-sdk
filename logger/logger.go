// Logging provider
//
// Log only fully init when app.Run() called
package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Fields logrus.Fields

var currentServLog ServiceLogger

func InitServLogger(allowFileLogger bool) {
	if allowFileLogger {
		currentServLog = DefaultMessageLogger
	} else {
		currentServLog = DefaultStdLogger
	}
}

// GetLogger returns a logger with the given name.
func GetCurrent() ServiceLogger {
	return currentServLog
}

type Logger interface {
	Print(args ...interface{})
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})

	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})

	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})

	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})

	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	Panic(...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})

	With(key string, value interface{}) Logger
	Withs(Fields) Logger
	// add source field to log
	WithSrc() Logger
	GetLevel() string
}

type logger struct {
	*logrus.Entry
}

// ServiceLogger is a logger for the service.
func (l *logger) GetLevel() string {
	return l.Entry.Logger.Level.String()
}

func (l *logger) debugSrc() *logrus.Entry {

	if _, ok := l.Entry.Data["source"]; ok {
		return l.Entry
	}

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	return l.Entry.WithField("source", fmt.Sprintf("%s:%d", file, line))
}

// Debug logs a message at level Debug on the standard logger.
func (l *logger) Debug(args ...interface{}) {
	if l.Entry.Logger.Level >= logrus.DebugLevel {
		l.debugSrc().Debug(args...)
	}
}

// Debugln logs a message at level Debug on the standard logger.
func (l *logger) Debugln(args ...interface{}) {
	if l.Entry.Logger.Level >= logrus.DebugLevel {
		l.debugSrc().Debugln(args...)
	}
}

// Debugf logs a message at level Debug on the standard logger.
func (l *logger) Debugf(format string, args ...interface{}) {
	if l.Entry.Logger.Level >= logrus.DebugLevel {
		l.debugSrc().Debugf(format, args...)
	}
}

// Print logs a message at level Debug on the standard logger.
func (l *logger) Print(args ...interface{}) {
	if l.Entry.Logger.Level >= logrus.DebugLevel {
		l.debugSrc().Debug(args...)
	}
}

// With returns a new logger with the given key and value.
func (l *logger) With(key string, value interface{}) Logger {
	return &logger{l.Entry.WithField(key, value)}
}

// Withs returns a new logger with the given fields.
func (l *logger) Withs(fields Fields) Logger {
	return &logger{l.Entry.WithFields(logrus.Fields(fields))}
}

// WithSrc returns a new logger with the source field.
func (l *logger) WithSrc() Logger {
	return &logger{l.debugSrc()}
}

// mustParseLevel parses a log level string and returns the corresponding logrus log level.
func mustParseLevel(level string) logrus.Level {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatal(err.Error())
	}
	return lv
}
