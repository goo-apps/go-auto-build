package goautobuild

import (
	"github.com/goo-apps/go-auto-build/watcher"
)

// Config is an alias to watcher.Config to expose it from root.
type Config = watcher.Config

// GoBuildWatcher is an alias to watcher.GoBuildWatcher for easy access.
type GoBuildWatcher = watcher.GoBuildWatcher

// NewWatcher creates a new build watcher.
var NewWatcher = watcher.NewWatcher
