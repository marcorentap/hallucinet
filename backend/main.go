package main

import (
	"log"

	"github.com/marcorentap/hallucinet/backend/client"
	"github.com/marcorentap/hallucinet/backend/committer"
	"github.com/marcorentap/hallucinet/backend/core"
	"github.com/marcorentap/hallucinet/backend/mapper"
)

func handleExistingContainers(hctx core.HallucinetContext) {
	client := hctx.Client
	mapper := hctx.Mapper
	for _, container := range client.GetContainers() {
		mapper.AddContainer(container)
	}
	hctx.Committer.Commit()
}

func handleContainerEvents(hctx core.HallucinetContext) {
	mapper := hctx.Mapper
	events := hctx.Client.GetEvents()

	for event := range events {
		switch event.EventType {
		case core.EventConnected:
			mapper.AddContainer(event.ContainerID)
		case core.EventDisconnected:
			mapper.RemoveContainer(event.ContainerID)
		}

		hctx.Committer.Commit()
	}
}

func main() {
	hctx := core.HallucinetContext{}
	hctx.Config = core.HallucinetConfig{
		Client:      "docker",
		Mapper:      "container_name",
		NetworkName: "hallucinet",
		Committer:   "hosts",
	}

	if hctx.Config.Client == "docker" {
		hctx.Client = client.NewDockerContainerClient(hctx)
	} else {
		log.Panicf("Unimplemented container client %v\n", hctx.Config.Client)
	}

	if hctx.Config.Mapper == "container_name" {
		hctx.Mapper = mapper.NewContainerNameMapper(hctx)
	} else {
		log.Panicf("Unimplemented mapper %v\n", hctx.Config.Mapper)
	}

	if hctx.Config.Committer == "hosts" {
		hctx.Committer = committer.NewHostsCommitter(hctx)
	} else {
		log.Panicf("Unimplemented committer %v\n", hctx.Config.Committer)
	}

	handleExistingContainers(hctx)
	handleContainerEvents(hctx)
}
