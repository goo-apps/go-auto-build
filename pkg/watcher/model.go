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
	ConfigPath    string `json:"config_path"`
	OutputBinary  string `json:"output_binary"`
	InstallPath   string `json:"install_path"`
	PollInterval  int    `json:"poll_interval_seconds"`
	WatchExt      string `json:"watch_ext"`
	ProjectRoot   string `json:"project_root"`
	EnableLogging bool   `json:"enable_logging"`
	BuildCommand  string `json:"build_command"`
	PostBuildMove bool   `json:"post_build_move"`
}

type GoBuildWatcher struct {
	Cfg        *Config
	ModTimeMap map[string]time.Time
	Mu         sync.Mutex
}
