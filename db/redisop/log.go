package redisop

import (
	"log"
	"os"
)

// Logger logger interface use by this lib
type Logger interface {
	Warning(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

const (
	levelWarn = iota + 1
	levelInfo
	levelDebug
)

func getLevel(level string) (l int) {
	switch level {
	case "warn":
		l = levelWarn
	case "info":
		l = levelInfo
	case "debug":
		l = levelDebug
	default:
		l = 0
	}
	return
}

// DefaultLogger implement default logger interface
type DefaultLogger struct {
	Level int
	l     *log.Logger
}

// NewDefaultLogger init new default logger
func NewDefaultLogger(level string) *DefaultLogger {
	ll := new(DefaultLogger)
	ll.Level = getLevel(level)
	ll.l = log.New(os.Stdout, "redisop ", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return ll
}

// Warning warning message
func (l DefaultLogger) Warning(format string, v ...interface{}) {
	if l.Level >= levelWarn {
		l.l.Printf(format, v...)
	}
	return
}

// Info info message
func (l DefaultLogger) Info(format string, v ...interface{}) {
	if l.Level >= levelInfo {
		l.l.Printf(format, v...)
	}
	return
}

// Debug debug message
func (l DefaultLogger) Debug(format string, v ...interface{}) {
	if l.Level >= levelDebug {
		l.l.Printf(format, v...)
	}
	return
}

var mylog Logger

// SetLogger set redisop info output
func SetLogger(info Logger) {
	mylog = info
}

func info(format string, v ...interface{}) {
	if mylog == nil {
		return
	}
	mylog.Info(format, v...)
}

func warn(format string, v ...interface{}) {
	if mylog == nil {
		return
	}
	mylog.Warning(format, v...)
}

func debug(format string, v ...interface{}) {
	if mylog == nil {
		return
	}
	mylog.Debug(format, v...)
}
