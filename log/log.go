package log

import (
	"fmt"
	syslog "log"
)

const (
	None Level = iota
	Fatal
	Error
	Warn
	Info
	Debug
)

var logging = struct {
	level  Level
	format string
}{
	level:  Warn,
	format: "%-5v  %-9s  %s",
}

type Level int

func (l Level) String() string {
	return [...]string{"", "FATAL", "ERROR", "WARN", "info", "debug"}[l]
}

func init() {
	syslog.SetFlags(0)
}

func SetLogLevel(l Level) {
	logging.level = l
}

func Debugf(operation string, format string, v ...interface{}) {
	log(Debug, operation, fmt.Sprintf(format, v...))
}

func Infof(operation string, format string, v ...interface{}) {
	log(Info, operation, fmt.Sprintf(format, v...))
}

func Warnf(operation string, format string, v ...interface{}) {
	log(Warn, operation, fmt.Sprintf(format, v...))
}

func Errorf(operation string, format string, v ...interface{}) {
	log(Error, operation, fmt.Sprintf(format, v...))
}

func Fatalf(operation string, format string, v ...interface{}) {
	syslog.Fatalf(logging.format, Fatal, operation, fmt.Sprintf(format, v...))
}

func log(l Level, operation string, msg string) {
	if l >= logging.level {
		syslog.Printf(logging.format, l, operation, msg)
	}
}
