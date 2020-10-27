package logger

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

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
		beegoLogConfig = fmt.Sprintf(`{"filename":%q,"maxsize":%d,"rotate":true,"maxdays":%d,"level":%d}`, logFileName, 1<<29, cfg.MaxDays, logs.LevelDebug)
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
