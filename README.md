# go-auto-build

A pluggable Go package build system using only Go standard libraries. This package can be imported into any static Go project to enable automatic builds after any code change during development.

## Features

- ⚡ Automatic rebuilds on file changes
- 🔌 Pluggable and easy to integrate
- 🛠️ Zero external dependencies—uses only Go standard libraries
- 🚀 Boosts productivity during development

## Installation

```sh
go get github.com/goo-apps/go-auto-build

## Use
cfg := &gobuildwatcher.Config{
		ConfigPath:     "resource.toml",
		OutputBinary:   "build/mybin",
		InstallPath:    "/usr/local/bin/mybin",
		PollInterval:   2,
		WatchExt:       ".go,.json",
		ProjectRoot:    ".",
		EnableLogging:  true,
		PostBuildMove:  true,
		BuildCommand:   "build -o build/mybin",
	}

	watcher := gobuildwatcher.NewWatcher(cfg)
	go watcher.Start() // runs in background
	select {}          // block forever
