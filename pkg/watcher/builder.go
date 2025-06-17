package watcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (w *GoBuildWatcher) buildAndInstall() error {
	cmdArgs := []string{"build", "-o", w.Cfg.OutputBinary}
	if w.Cfg.BuildCommand != "" {
		cmdArgs = strings.Fields(w.Cfg.BuildCommand)
	}
	cmd := exec.Command("go", cmdArgs...)
	cmd.Env = append(os.Environ(), "CONFIG_PATH="+w.Cfg.ConfigPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build error: %w", err)
	}

	if w.Cfg.PostBuildMove {
		mvCmd := exec.Command("sudo", "mv", w.Cfg.OutputBinary, w.Cfg.InstallPath)
		mvCmd.Stdout = os.Stdout
		mvCmd.Stderr = os.Stderr
		if err := mvCmd.Run(); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}
	}
	return nil
}
