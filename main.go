package main

import (
	"ail/asdf"
	"ail/updater"
	"log"
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
			updater.Update(p)
			wg.Done()
		}(plugin)
	}
	wg.Wait()

	err = asdf.Reshim()
	if err != nil {
		log.Fatalf("failed to reshim: %s", err)
	}
}
