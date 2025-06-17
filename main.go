package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goo-apps/go-auto-build/logger"
	"github.com/goo-apps/go-auto-build/pkg/watcher"
)

// CLI Entrypoint
func main() {
	cfgPath := flag.String("config", "watcher.json", "Path to config file (JSON|TOML|ENV)")
	once := flag.Bool("once", false, "Run build once and exit")
	root := flag.String("root", ".", "Root directory to watch")
	flag.Parse()

	cfg := &watcher.Config{}
	if err := loadConfig(*cfgPath, cfg); err != nil {
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

// Config loading
func loadConfig(path string, cfg *watcher.Config) error {
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		return loadJSONConfig(path, cfg)
	case ".env":
		return loadEnvConfig(path, cfg)
	case ".toml":
		return loadTOMLConfig(path, cfg)
	default:
		return fmt.Errorf("unsupported config format: %s", ext)
	}
}

func loadJSONConfig(path string, cfg *watcher.Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(cfg)
}

func loadEnvConfig(path string, cfg *watcher.Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	return loadConfigFromEnv(cfg)
}

func loadTOMLConfig(path string, cfg *watcher.Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	configMap := make(map[string]string)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		configMap[key] = val
	}
	applyMapToConfig(configMap, cfg)
	return nil
}

func loadConfigFromEnv(cfg *watcher.Config) error {
	cfg.ConfigPath = os.Getenv("CONFIG_PATH")
	cfg.OutputBinary = os.Getenv("OUTPUT_BINARY")
	cfg.InstallPath = os.Getenv("INSTALL_PATH")
	cfg.ProjectRoot = os.Getenv("PROJECT_ROOT")
	cfg.WatchExt = os.Getenv("WATCH_EXT")
	cfg.BuildCommand = os.Getenv("BUILD_COMMAND")
	cfg.EnableLogging = os.Getenv("ENABLE_LOGGING") == "true"
	cfg.PostBuildMove = os.Getenv("POST_BUILD_MOVE") == "true"
	if interval := os.Getenv("POLL_INTERVAL"); interval != "" {
		fmt.Sscanf(interval, "%d", &cfg.PollInterval)
	}
	return nil
}

func applyMapToConfig(m map[string]string, cfg *watcher.Config) {
	cfg.ConfigPath = m["config_path"]
	cfg.OutputBinary = m["output_binary"]
	cfg.InstallPath = m["install_path"]
	cfg.ProjectRoot = m["project_root"]
	cfg.WatchExt = m["watch_ext"]
	cfg.BuildCommand = m["build_command"]
	cfg.EnableLogging = m["enable_logging"] == "true"
	cfg.PostBuildMove = m["post_build_move"] == "true"
	fmt.Sscanf(m["poll_interval_seconds"], "%d", &cfg.PollInterval)
}
