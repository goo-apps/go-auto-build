// Author: rohan.das

// go-auto-builder - A pluggable go package build with go standard libraries, which can be imported to any static project for auto build after any change of the code during development. 
// Copyright (c) 2025 Go Application Hub @Rohan (rohan.das1203@gmail.com)
// Licensed under the MIT License. See LICENSE file for details.

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
