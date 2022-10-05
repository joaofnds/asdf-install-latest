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
}

func Update(pkg string) {
	l := NewPackageLogger(pkg)
	if config.Ignore[pkg] {
		l.Println("found it on ignore list, skipping...")
		return
	}

	l.Println("started", pkg)
	curr, err := asdf.Current(pkg)
	if err != nil {
		l.Println("failed to get current verison")
		return
	}

	latest, err := asdf.Latest(pkg)
	if err != nil {
		l.Println("failed to get latest verison")
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

	err = runHook(pkg)
	if err != nil {
		l.Printf("failed to run hook: %s\n", err)
		return
	}
}

func runHook(pkg string) error {
	hook := path.Join(config.HooksDir, pkg+".sh")
	info, err := os.Stat(hook)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}

	if !isExecutable(info) {
		return fmt.Errorf("hook is not executable, please run: chmod +x %q", hook)
	}

	cmd := &exec.Cmd{
		Path:   hook,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}

	err = cmd.Run()
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
