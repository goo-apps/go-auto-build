package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goo-apps/go-auto-build/watcher"
	"gopkg.in/yaml.v3"
)

// Config loading
func LoadConfig(path string, cfg *watcher.Config) error {
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		return loadJSONConfig(path, cfg)
	case ".toml":
		return loadTOMLConfig(path, cfg)
	case ".yaml", ".yml":
		return loadYAMLConfig(path, cfg)
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

func loadYAMLConfig(path string, cfg *watcher.Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	return decoder.Decode(cfg)
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
