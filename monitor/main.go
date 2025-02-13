package main

import (
	"fmt"

	"github.com/marcorentap/hallucinet/config"
	"github.com/marcorentap/hallucinet/types"
	"github.com/marcorentap/hallucinet/watcher"
)

func main() {
	config := config.NewHallucinetConfig()
	hctx := types.HallucinetContext{
		Config:    config,
		EventChan: make(chan types.HallucinetEvent),
	}

	go watcher.WatchDockerEvents(hctx)
	for event := range hctx.EventChan {
		fmt.Println(event)
	}
}
