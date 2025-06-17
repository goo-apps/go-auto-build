package watcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)


func (w *GoBuildWatcher) buildAndInstall() error {
	cmdArgs := []string{"build", "-o", w.cfg.OutputBinary}
	if w.cfg.BuildCommand != "" {
		cmdArgs = strings.Fields(w.cfg.BuildCommand)
	}
	cmd := exec.Command("go", cmdArgs...)
	cmd.Env = append(os.Environ(), "CONFIG_PATH="+w.cfg.ConfigPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build error: %w", err)
	}

	if w.cfg.PostBuildMove {
		mvCmd := exec.Command("sudo", "mv", w.cfg.OutputBinary, w.cfg.InstallPath)
		mvCmd.Stdout = os.Stdout
		mvCmd.Stderr = os.Stderr
		if err := mvCmd.Run(); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}
	}
	return nil
}

