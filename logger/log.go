package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/beego/beego/v2/core/logs"
)

var moreLogger sync.Map

// Logger logger interface use by this lib
type Logger interface {
	Warning(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type logObj struct {
	bl *logs.BeeLogger
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if !strings.Contains(msg, "%") {
			// do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

func (l *logObj) Critical(format string, v ...interface{}) {
	l.bl.Critical(formatLog(format, v...))
}

func (l *logObj) Error(format string, v ...interface{}) {
	l.bl.Error(formatLog(format, v...))
}

func (l *logObj) Warning(format string, v ...interface{}) {
	l.bl.Warning(formatLog(format, v...))
}

func (l *logObj) Info(format string, v ...interface{}) {
	l.bl.Info(formatLog(format, v...))
}

func (l *logObj) Debug(format string, v ...interface{}) {
	l.bl.Debug(formatLog(format, v...))
}

var lo logObj = logObj{
	bl: logs.GetBeeLogger(),
}

// GetLogger get logger interface
func GetLogger() Logger {
	return &lo
}

type logSetupConfig struct {
	MaxDays  int    `json:"maxdays"`
	ToSyslog bool   `json:"tosyslog,omitempty"`
	AppName  string `json:"appname,omitempty"`
	FileName string `json:"filename,omitempty"`
}

// Setup init log settings, level fatal,error,warn,info,debug
func Setup(console bool, level string, logFileName string, jsonConfig ...string) {
	var cfg logSetupConfig
	if len(jsonConfig) > 0 && jsonConfig[0] != "" {
		json.Unmarshal([]byte(jsonConfig[0]), &cfg)
	}
	if logFileName != "" {
		var beegoLogConfig string
		if cfg.MaxDays <= 0 {
			cfg.MaxDays = 7
		}
		beegoLogConfig = fmt.Sprintf(`{"filename":%q,"maxsize":%d,"useopentime":true,"rotate":true,"maxdays":%d,"level":%d}`, logFileName, 1<<29, cfg.MaxDays, logs.LevelDebug)
		logs.SetLogger("multifile", beegoLogConfig)
	}
	if cfg.ToSyslog == true && cfg.AppName != "" {
		logs.SetLogger("syslog", fmt.Sprintf(`{"appname":%q,"level":%d}`, cfg.AppName, logs.LevelDebug))
	}
	if console {
		logs.SetLogger("console")
	} else {
		logs.GetBeeLogger().DelLogger("console")
	}
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(4)
	logs.SetLevel(getLevel(level))
	return
}

// SetLevel config log level
func SetLevel(level string) {
	logs.SetLevel(getLevel(level))
}

// Fatal logs a message at critical level
func Fatal(f interface{}, v ...interface{}) {
	logs.Critical(f, v...)
}

// Error logs a message at error level
func Error(f interface{}, v ...interface{}) {
	logs.Error(f, v...)
}

// Warning logs a message at warning level
func Warning(f interface{}, v ...interface{}) {
	logs.Warning(f, v...)
}

// Info logs a message at info level
func Info(f interface{}, v ...interface{}) {
	logs.Informational(f, v...)
}

// Debug logs a message at debug level
func Debug(f interface{}, v ...interface{}) {
	logs.Debug(f, v...)
}

// NewLogFile init another log file
func NewLogFile(logFileName string, maxDays int, tags ...string) (newLog *logs.BeeLogger) {
	var beegoLogConfig string
	newLog = logs.NewLogger()
	beegoLogConfig = fmt.Sprintf(`{"filename":%q,"maxsize":%d,"useopentime":true,"rotate":true,"maxdays":%d,"level":%d}`, logFileName, 1<<29, maxDays, logs.LevelDebug)
	newLog.SetLogger("multifile", beegoLogConfig)
	newLog.DelLogger("console")
	newLog.EnableFuncCallDepth(false)
	newLog.SetLogFuncCallDepth(3)
	newLog.SetLevel(getLevel("info"))
	if len(tags) > 0 {
		moreLogger.Store(tags[0], newLog)
	}
	return
}

func getBeeLog(tag string) (log *logs.BeeLogger) {
	if value, ok := moreLogger.Load(tag); ok == true {
		log, ok = value.(*logs.BeeLogger)
	}
	return
}

// EnalbeFuncCallFor enable log call func
func EnalbeFuncCallFor(tag string, enabled ...bool) {
	if ll := getBeeLog(tag); ll != nil {
		var toEnable = true
		if len(enabled) > 0 {
			toEnable = enabled[0]
		}
		ll.EnableFuncCallDepth(toEnable)
	}
}

// EnalbeConsoleFor enable log to console
func EnalbeConsoleFor(tag string, enabled ...bool) {
	if ll := getBeeLog(tag); ll != nil {
		var toEnable = true
		if len(enabled) > 0 {
			toEnable = enabled[0]
		}
		if toEnable {
			ll.SetLogger("console")
		} else {
			ll.DelLogger("console")
		}
	}
}

// ErrorFor logs a message at error level
func ErrorFor(tag string, f interface{}, v ...interface{}) {
	if ll := getBeeLog(tag); ll != nil {
		ll.Error(formatLog(f, v...))
	}
}

// WarningFor logs a message at warning level
func WarningFor(tag string, f interface{}, v ...interface{}) {
	if ll := getBeeLog(tag); ll != nil {
		ll.Warning(formatLog(f, v...))
	}
}

// InfoFor logs a message at info level
func InfoFor(tag string, f interface{}, v ...interface{}) {
	if ll := getBeeLog(tag); ll != nil {
		ll.Informational(formatLog(f, v...))
	}
}

// getLevel fatal,error,warn,info,debug
func getLevel(level string) (l int) {
	switch level {
	case "fatal":
		l = logs.LevelCritical
	case "error":
		l = logs.LevelError
	case "warn":
		l = logs.LevelWarning
	case "info":
		l = logs.LevelInformational
	case "debug":
		l = logs.LevelDebug
	default:
		l = logs.LevelInformational
	}
	return
}
