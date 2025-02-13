package watcher

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/marcorentap/hallucinet/types"
)

func createDockerClient() *client.Client {
	cli, cliErr := client.NewClientWithOpts()
	if cliErr != nil {
		log.Panicf("Cannot create docker client: %v\n", cliErr)
	}
	return cli
}

func createDockerChannel(cli *client.Client) (<-chan events.Message, <-chan error) {
	filters := filters.NewArgs(
		filters.Arg("type", "network"),
		filters.Arg("event", "connect"),
		filters.Arg("event", "disconnect"),
	)

	dockerChan, chanErr := cli.Events(context.Background(),
		events.ListOptions{Filters: filters})
	return dockerChan, chanErr
}

func getContainerName(cli *client.Client, containerID string) string {
	conJson, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Panicf("Cannot inspect container %v: %v\n", containerID, err)
	}
	return conJson.Name
}

func TranslateDockerEvent(cli *client.Client, e events.Message) types.HallucinetEvent {
	time := time.Now()
	networkID := e.Actor.ID
	networkName := e.Actor.Attributes["name"]
	containerID := e.Actor.Attributes["container"]
	containerName := getContainerName(cli, containerID)[1:]

	var kind types.HallucinetEventKind
	switch e.Action {
	case events.ActionConnect:
		kind = types.ContainerConnected
	case events.ActionDisconnect:
		kind = types.ContainerDisconnected
	default:
		kind = types.UnknownEvent
	}

	return types.HallucinetEvent{
		Kind:          kind,
		ContainerID:   containerID,
		ContainerName: containerName,
		NetworkID:     networkID,
		NetworkName:   networkName,
		Data:          nil,
		Time:          time,
	}
}

func publishExistingContainers(cli *client.Client, eventChan chan types.HallucinetEvent) {
	networks, networkErr := cli.NetworkList(context.Background(), network.ListOptions{})
	if networkErr != nil {
		log.Panicf("Cannot obtain docker networks: %v\n", networkErr)
	}
	networkIDToName := make(map[string]string, len(networks))
	for _, network := range networks {
		networkID := network.ID
		networkName := network.Name
		networkIDToName[networkID] = networkName
	}

	containers, containerErr := cli.ContainerList(context.Background(), container.ListOptions{})
	if containerErr != nil {
		log.Printf("Cannot obtain docker containers: %v\n", containerErr)
	}
	for _, container := range containers {
		containerID := container.ID
		containerName := container.Names[0][1:]

		for networkName, setting := range container.NetworkSettings.Networks {
			networkID := setting.NetworkID
			eventChan <- types.HallucinetEvent{
				Kind:          types.ContainerConnected,
				ContainerID:   containerID,
				ContainerName: containerName,
				NetworkID:     networkID,
				NetworkName:   networkName,
				Data:          nil,
				Time:          time.Now(),
			}
		}
	}
}

func WatchDockerEvents(hctx types.HallucinetContext) {
	cli := createDockerClient()
	defer cli.Close()

	dockerChan, dockerErrChan := createDockerChannel(cli)
	publishExistingContainers(cli, hctx.EventChan)
	for {
		select {
		case dockerEvent := <-dockerChan:
			event := TranslateDockerEvent(cli, dockerEvent)
			hctx.EventChan <- event
		case dockerErr := <-dockerErrChan:
			log.Panicf("Docker error event: %v\n", dockerErr)
		}
	}
}
