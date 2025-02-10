package main

import (
	"log"
	"os"

	"github.com/marcorentap/hallucinet/client"
	"github.com/marcorentap/hallucinet/committer"
	"github.com/marcorentap/hallucinet/core"
	"github.com/marcorentap/hallucinet/mapper"
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

func getEnvOrDefault(envar string, def string) string {
	val, present := os.LookupEnv(envar)
	if !present {
		return def
	}
	return val
}

func main() {
	hctx := core.HallucinetContext{}
	hctx.Config = core.HallucinetConfig{
		Client:      getEnvOrDefault("CLIENT", "docker"),
		Mapper:      getEnvOrDefault("MAPPER", "container_name"),
		NetworkName: getEnvOrDefault("NETWORK", "hallucinet"),
		Committer:   getEnvOrDefault("COMMITTER", "hosts"),
		Suffix:      getEnvOrDefault("SUFFIX", ".test"),
		HostsPath:   "/var/hallucinet/hosts",
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
