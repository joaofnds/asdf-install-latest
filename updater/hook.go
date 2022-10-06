package updater

import (
	"ail/config"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

func getHook(pkg string) (*exec.Cmd, error) {
	hook := config.HookPath(pkg)

	info, err := os.Stat(hook)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if !isExecutable(info) {
		return nil, fmt.Errorf("hook is not executable, please run: chmod +x %q", hook)
	}

	return exec.Command(config.HookPath(pkg)), nil
}

func isExecutable(info fs.FileInfo) bool {
	return info.Mode()&0100 != 0
}
