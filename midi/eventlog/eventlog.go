package eventlog

import (
	"fmt"
	"os"
)

var EventLog = struct {
	Verbose bool
	Debug   bool
}{
	Verbose: false,
	Debug:   false,
}

func Debug(msg string) {
	if EventLog.Debug {
		fmt.Fprintf(os.Stdout, "(debug) %s\n", msg)
	}
}

func Info(msg string) {
	if EventLog.Verbose {
		fmt.Fprintf(os.Stdout, "(info)  %s\n", msg)
	}
}

func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "WARN  %s\n", msg)
}

func Error(msg string) {
	fmt.Fprintf(os.Stderr, "ERROR %s\n", msg)
}
