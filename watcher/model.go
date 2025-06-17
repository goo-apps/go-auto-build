// Author: rohan.das

// go-auto-builder - A pluggable go package build with go standard libraries, which can be imported to any static project for auto build after any change of the code during development.
// Copyright (c) 2025 Go Application Hub @Rohan (rohan.das1203@gmail.com)
// Licensed under the MIT License. See LICENSE file for details.

package watcher

import (
	"sync"
	"time"
)

type Watcher interface {
	Start()
	RunOnce()
}

type Config struct {
	ConfigPath    string   `json:"config_path"`
	OutputBinary  string   `json:"output_binary"`
	InstallPath   string   `json:"install_path"`
	PollInterval  int      `json:"poll_interval_seconds"`
	WatchExt      string   `json:"watch_ext"`
	ProjectRoot   string   `json:"project_root"`
	EnableLogging bool     `json:"enable_logging"`
	BuildCommand  string   `json:"build_command"`
	PostBuildMove bool     `json:"post_build_move"`
	ExcludePaths  []string `json:"exclude_paths" yaml:"exclude_paths"`
}

type GoBuildWatcher struct {
	cfg        *Config
	modTimeMap map[string]time.Time
	mu         sync.Mutex
}
