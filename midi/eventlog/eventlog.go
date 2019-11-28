package eventlog

import (
	"fmt"
	"os"
)

func Debug(msg string) {
	fmt.Fprintf(os.Stdout, "DEBUG: %s\n", msg)
}

func Info(msg string) {
	fmt.Fprintf(os.Stdout, "INFO:  %s\n", msg)
}

func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "WARN:  %s\n", msg)
}

func Error(msg string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", msg)
}
