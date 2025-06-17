package watcher

import (
	"os"
	"path/filepath"
	"time"

	"github.com/goo-apps/go-auto-build/logger"
)


func NewWatcher(cfg *Config) *GoBuildWatcher {
	return &GoBuildWatcher{
		Cfg:        cfg,
		ModTimeMap: make(map[string]time.Time),
	}
}

func (w *GoBuildWatcher) Start() {
	logger.Info("Watching %s for *%s file changes every %ds...", w.Cfg.ProjectRoot, w.Cfg.WatchExt, w.Cfg.PollInterval)
	for {
		w.RunOnce()
		time.Sleep(time.Duration(w.Cfg.PollInterval) * time.Second)
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
	w.Mu.Lock()
	defer w.Mu.Unlock()

	var changed bool
	err := filepath.WalkDir(w.Cfg.ProjectRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != w.Cfg.WatchExt {
			return nil
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil
		}
		last, seen := w.ModTimeMap[path]
		if !seen || info.ModTime().After(last) {
			w.ModTimeMap[path] = info.ModTime()
			changed = true
		}
		return nil
	})
	return changed, err
}
