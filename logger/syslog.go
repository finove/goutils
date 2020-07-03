// +build !windows,!plan9

package logger

import (
	"encoding/json"
	"log/syslog"
	"time"

	"github.com/astaxie/beego/logs"
)

// adapter name
const (
	AdapterSyslog = "syslog"
)

type syslogWriter struct {
	innerWriter *syslog.Writer
	AppName     string `json:"appname"`
	Level       int    `json:"level"`
}

// NewSyslog create new syslogWriter returning as LoggerInterface.
func NewSyslog() logs.Logger {
	conn := new(syslogWriter)
	conn.Level = logs.LevelTrace
	return conn
}

// Init init connection writer with json config.
// json config only need key "level".
func (c *syslogWriter) Init(jsonConfig string) (err error) {
	err = json.Unmarshal([]byte(jsonConfig), c)
	if err != nil {
		return
	}
	c.innerWriter, err = syslog.New(syslog.LOG_DEBUG|syslog.LOG_LOCAL1, c.AppName)
	if err != nil {
		return
	}
	return
}

// WriteMsg write message in syslog.
func (c *syslogWriter) WriteMsg(when time.Time, msg string, level int) (err error) {
	if level > c.Level {
		return nil
	}
	if c.innerWriter == nil {
		c.innerWriter, err = syslog.New(syslog.LOG_DEBUG|syslog.LOG_LOCAL1, c.AppName)
		if err != nil {
			return
		}
	}
	switch level {
	case logs.LevelEmergency:
		c.innerWriter.Emerg(msg)
	case logs.LevelAlert:
		c.innerWriter.Alert(msg)
	case logs.LevelCritical:
		c.innerWriter.Crit(msg)
	case logs.LevelError:
		c.innerWriter.Err(msg)
	case logs.LevelWarning:
		c.innerWriter.Warning(msg)
	case logs.LevelNotice:
		c.innerWriter.Notice(msg)
	case logs.LevelInformational:
		c.innerWriter.Info(msg)
	case logs.LevelDebug:
		c.innerWriter.Debug(msg)
	default:
		c.innerWriter.Debug(msg)
	}
	return nil
}

// Flush implementing method. empty.
func (c *syslogWriter) Flush() {
}

// Destroy destroy syslog writer
func (c *syslogWriter) Destroy() {
	if c.innerWriter != nil {
		c.innerWriter.Close()
		c.innerWriter = nil
	}
}

func init() {
	logs.Register(AdapterSyslog, NewSyslog)
}
