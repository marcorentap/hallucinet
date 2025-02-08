package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/marcorentap/hallucinet/backend/internal/mapper"
)

type HallucinateConfig struct {
	networkName string
	mapper      mapper.Mapper
}

var config HallucinateConfig

func handleExistingContainers(cli *client.Client, ctx context.Context) {
	networkName := config.networkName
	mapper := config.mapper
	ins, err := cli.NetworkInspect(ctx, networkName, network.InspectOptions{})
	if err != nil {
		log.Fatalf("Cannot inspect network %v: %v", networkName, err)
	}

	for containerID := range ins.Containers {
		mapper.AddContainer(containerID)
	}
}

func handleConnectionEvents(cli *client.Client, ctx context.Context) {
	mapper := config.mapper
	eventFilters := filters.NewArgs(
		filters.Arg("type", "network"),
		filters.Arg("event", "connect"),
		filters.Arg("event", "disconnect"),
	)
	msgChan, errChan := cli.Events(ctx, events.ListOptions{Filters: eventFilters})

	for {
		select {
		case msg := <-msgChan:
			attr := msg.Actor.Attributes
			action := msg.Action

			networkName := attr["name"]
			if networkName != config.networkName {
				continue
			}

			containerID := attr["container"]
			if action == events.ActionConnect {
				mapper.AddContainer(containerID)
			} else if action == events.ActionDisconnect {
				mapper.RemoveContainer(containerID)
			}

			mapper.ToHosts()

		case err := <-errChan:
			log.Panicf("Cannot handle Docker events: %v\n", err)
		}
	}
}

func createDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panicf("Cannot create Docker client: %v\n", err)
	}
	return cli
}

func main() {
	ctx := context.Background()
	cli := createDockerClient()
	defer cli.Close()

	config = HallucinateConfig{
		networkName: "hallucinet",
		mapper:      mapper.NewContainerNameMapper(ctx, cli),
	}

	handleExistingContainers(cli, ctx)
	handleConnectionEvents(cli, ctx)
}
