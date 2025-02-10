package client

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/marcorentap/hallucinet/core"
)

type DockerContainerClient struct {
	hctx core.HallucinetContext
	cli  *client.Client
}

func NewDockerContainerClient(hctx core.HallucinetContext) *DockerContainerClient {
	c := DockerContainerClient{hctx: hctx}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panicf("Cannot create Docker client: %v\n", err)
	}
	c.cli = cli
	return &c
}

func (c *DockerContainerClient) GetContainers() []string {
	networkName := c.hctx.Config.NetworkName
	inspect, err := c.cli.NetworkInspect(context.Background(), networkName, network.InspectOptions{})
	if err != nil {
		log.Panicf("Cannot inspect network %v: %v\n", networkName, err)
	}

	keys := make([]string, 0, len(inspect.Containers))
	for k := range inspect.Containers {
		keys = append(keys, k)
	}
	return keys
}

func (c *DockerContainerClient) GetContainerName(containerID string) string {
	con, err := c.cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Panicf("Cannot inspect container %v: %v\n", containerID, err)
	}
	return con.Name
}

func (c *DockerContainerClient) GetContainerAddr(containerID string) string {
	con, err := c.cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Panicf("Cannot inspect container %v: %v\n", containerID, err)
	}

	networks := con.NetworkSettings.Networks
	endpoint, exists := networks[c.hctx.Config.NetworkName]

	if !exists {
		return ""
	}

	return endpoint.IPAddress
}

func (c *DockerContainerClient) GetEvents() <-chan core.ContainerEvent {
	eventChan := make(chan core.ContainerEvent)
	eventFilters := filters.NewArgs(
		filters.Arg("type", "network"),
		filters.Arg("event", "connect"),
		filters.Arg("event", "disconnect"),
	)

	ctx := context.Background()
	msgChan, errChan := c.cli.Events(ctx, events.ListOptions{Filters: eventFilters})

	go func() {
		for {
			select {
			case msg := <-msgChan:
				attr := msg.Actor.Attributes
				networkName := attr["name"]
				containerID := attr["container"]
				if networkName != c.hctx.Config.NetworkName {
					continue
				}

				var eventType core.ContainerEventType
				if msg.Action == events.ActionConnect {
					eventType = core.EventConnected
				} else if msg.Action == events.ActionDisconnect {
					eventType = core.EventDisconnected
				} else {
					eventType = core.EventUnknown
				}

				eventChan <- core.ContainerEvent{
					ContainerID: containerID,
					EventType:   eventType,
				}

			case err := <-errChan:
				log.Panicf("Error reading from docker events: %v\n", err)
			}

		}
	}()

	return eventChan
}
