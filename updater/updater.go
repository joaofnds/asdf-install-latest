package updater

import (
	"ail/asdf"
	"ail/config"
	"log"
	"os"
)

func Update(pkg string) {
	logger := log.New(os.Stdout, "["+pkg+"] ", log.Lmsgprefix)

	if config.IsIgnored(pkg) {
		logger.Println("ignoring...")
		return
	}

	curr, err := asdf.LatestInstalled(pkg)
	if err != nil {
		logger.Printf("failed to get latest version: %s\n", err)
		return
	}

	latest, err := asdf.Latest(pkg)
	if err != nil {
		logger.Printf("failed to get latest version: %s\n", err)
		return
	}

	if !curr.LessThan(latest) {
		logger.Println("no updates")
		return
	}

	logger.Printf("will update from %q to %q\n", curr, latest)

	err = asdf.Install(pkg, latest)
	if err != nil {
		logger.Printf("failed to install version %q\n", latest)
		return
	}

	logger.Printf("setting global version to %q\n", latest)

	err = asdf.SetGlobal(pkg, latest)
	if err != nil {
		logger.Printf("failed to set global version to %q\n", latest)
		return
	}

	hook, err := getHook(pkg)
	if err != nil {
		logger.Printf("hook error: %s", err)
	}

	if hook == nil {
		return
	}

	logger.Println("running hook")

	if err = hook.Run(); err != nil {
		logger.Printf("hook failed: %s\n", err)
	}
}
