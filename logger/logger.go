package logger

import (
	"fmt"
	"os"
)

// Logging helpers
func Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func Warn(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

func Fatal(format string, args ...interface{}) {
	fmt.Printf("[FATAL] "+format+"\n", args...)
	os.Exit(1)
}
