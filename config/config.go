package config

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path"
)

var (
	Ignore     = map[string]bool{}
	ConfigDir  string
	HooksDir   string
	IgnoreFile string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get user home dir: %s", err)
	}
	ConfigDir = path.Join(home, ".config", "ail")
	HooksDir = path.Join(ConfigDir, "hooks")
	IgnoreFile = path.Join(ConfigDir, "ignore")

	f, err := os.Open(IgnoreFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("could not open ignore file: %s", err)
	}

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanWords)
	for s.Scan() {
		Ignore[s.Text()] = true
	}
}

func HookPath(pkg string) string {
	return path.Join(HooksDir, pkg+".sh")
}
