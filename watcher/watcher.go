// Author: rohan.das

// go-auto-builder - A pluggable go package build with go standard libraries, which can be imported to any static project for auto build after any change of the code during development.
// Copyright (c) 2025 Go Application Hub @Rohan (rohan.das1203@gmail.com)
// Licensed under the MIT License. See LICENSE file for details.

package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goo-apps/go-auto-build/internal/logger"
)

func NewWatcher(cfg *Config) *GoBuildWatcher {
	return &GoBuildWatcher{
		cfg:        cfg,
		modTimeMap: make(map[string]time.Time),
	}
}

func (w *GoBuildWatcher) Start() {
	logger.Info("Watching %s for file changes every %ds...", w.cfg.ProjectRoot, w.cfg.PollInterval)
	for {
		w.RunOnce()
		time.Sleep(time.Duration(w.cfg.PollInterval) * time.Second)
	}
}

func (w *GoBuildWatcher) RunOnce() {
	changed, err := w.checkChanges()
	if err != nil {
		logger.Warn("Error checking changes: %v", err)
		return
	}
	if changed {
		logger.Info("Change detected. Rebuilding...")
		if err := w.buildAndInstall(); err != nil {
			logger.Warn("Build failed: %v", err)
		} else {
			logger.Info("Build and install complete.")
		}
	}
}

func (w *GoBuildWatcher) checkChanges() (bool, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var changed bool
	extFilters := strings.Split(w.cfg.WatchExt, ",")

	outputAbs, _ := filepath.Abs(w.cfg.OutputBinary)
	configAbs, _ := filepath.Abs(w.cfg.ConfigPath) // 🔧 <- get abs path to config file

	err := filepath.WalkDir(w.cfg.ProjectRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		fileAbs, _ := filepath.Abs(path)

		// ✅ Exclude binary and config file
		if fileAbs == outputAbs ||
			fileAbs == configAbs ||
			strings.Contains(path, ".git") ||
			strings.HasSuffix(path, ".DS_Store") ||
			w.isExcluded(fileAbs) {
			return nil
		}

		if len(extFilters) > 0 {
			match := false
			for _, ext := range extFilters {
				if strings.HasSuffix(path, strings.TrimSpace(ext)) {
					match = true
					break
				}
			}
			if !match {
				return nil
			}
		}

		info, err := os.Stat(path)
		if err != nil {
			return nil
		}
		last, seen := w.modTimeMap[path]
		if !seen || info.ModTime().After(last) {
			logger.Info("Detected change in: %s at %v)", path, info.ModTime())
			w.modTimeMap[path] = info.ModTime()
			changed = true
		}
		return nil
	})
	return changed, err
}

func (w *GoBuildWatcher) isExcluded(path string) bool {
	for _, excl := range w.cfg.ExcludePaths {
		exclAbs, err := filepath.Abs(excl)
		if err != nil {
			continue
		}
		if strings.HasPrefix(path, exclAbs) {
			return true
		}
	}
	return false
}
