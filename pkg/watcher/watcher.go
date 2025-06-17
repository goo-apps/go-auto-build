package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goo-apps/go-auto-build/logger"
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

	err := filepath.WalkDir(w.cfg.ProjectRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		fileAbs, _ := filepath.Abs(path)
		if fileAbs == outputAbs || strings.Contains(path, ".git") || strings.HasSuffix(path, ".DS_Store") {
			return nil
		}

		// Check extension filter
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
			w.modTimeMap[path] = info.ModTime()
			changed = true
		}
		return nil
	})
	return changed, err
}
