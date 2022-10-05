package main

import (
	"ail/asdf"
	"ail/config"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
)

func main() {
	plugins, err := asdf.PluginList()
	if err != nil {
		log.Fatalf("failed to get plugin list: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(plugins))
	for _, plugin := range plugins {
		go func(p string) {
			Update(p)
			wg.Done()
		}(plugin)
	}
	wg.Wait()

	err = asdf.Reshim()
	if err != nil {
		log.Fatalf("failed to reshim: %s", err)
	}
}

func Update(pkg string) {
	l := NewPackageLogger(pkg)
	if config.Ignore[pkg] {
		l.Println("found it on ignore list, skipping...")
		return
	}

	curr, err := asdf.Current(pkg)
	if err != nil {
		l.Printf("failed to get current verison: %s\n", err)
		return
	}

	latest, err := asdf.Latest(pkg)
	if err != nil {
		l.Printf("failed to get latest verison: %s\n", err)
		return
	}

	if !curr.LessThan(*latest) {
		l.Println("nothing to update")
		return
	}

	l.Printf("new version detected, will upgrade %q -> %q\n", curr, latest)

	err = asdf.Install(pkg, latest)
	if err != nil {
		l.Printf("failed to install verison %q\n", latest)
		return
	}

	l.Printf("setting global version of %q to %q\n", pkg, latest)
	err = asdf.SetGlobal(pkg, latest)
	if err != nil {
		l.Printf("failed to set global version to %q\n", latest)
		return
	}

	has, err := hasHook(pkg)
	if err != nil {
		l.Printf("a problem was found with the %q hook: %s", pkg, err)
	}

	if !has {
		return
	}

	l.Printf("running %q post install hook", pkg)

	err = runHook(pkg)
	if err != nil {
		l.Printf("failed to run hook: %s\n", err)
	}
}

func hasHook(pkg string) (bool, error) {
	hook := hookPath(pkg)
	info, err := os.Stat(hook)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return true, err
	}

	if !isExecutable(info) {
		return true, fmt.Errorf("hook is not executable, please run: chmod +x %q", hook)
	}

	return true, nil
}

func runHook(pkg string) error {
	err := exec.Command(hookPath(pkg)).Run()
	if err != nil {
		return fmt.Errorf("hook failed: %w", err)
	}

	return nil
}

func isExecutable(info fs.FileInfo) bool {
	return info.Mode()&0100 != 0
}

func NewPackageLogger(pkg string) log.Logger {
	return *log.New(os.Stdout, "["+pkg+"] ", log.Lmsgprefix)
}

func hookPath(pkg string) string {
	return path.Join(config.HooksDir, pkg+".sh")
}
