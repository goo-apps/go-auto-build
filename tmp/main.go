// Author: rohan.das

// go-auto-builder - A pluggable go package build with go standard libraries, which can be imported to any static project for auto build after any change of the code during development.
// Copyright (c) 2025 Go Application Hub @Rohan (rohan.das1203@gmail.com)
// Licensed under the MIT License. See LICENSE file for details.

package main

import (
	"flag"

	"github.com/goo-apps/go-auto-build/internal/config"
	"github.com/goo-apps/go-auto-build/internal/logger"
	"github.com/goo-apps/go-auto-build/watcher"
)

// CLI Entrypoint
func main() {
	cfgPath := flag.String("config", "watcher.json", "Path to config file (JSON|TOML|ENV)")
	once := flag.Bool("once", false, "Run build once and exit")
	root := flag.String("root", ".", "Root directory to watch")
	flag.Parse()

	cfg := &watcher.Config{}
	if err := config.LoadConfig(*cfgPath, cfg); err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}
	if cfg.ProjectRoot == "" {
		cfg.ProjectRoot = *root
	}

	watcher := watcher.NewWatcher(cfg)
	if *once {
		watcher.RunOnce()
		return
	}
	watcher.Start()
}
