package main

import (
	"context"
	"fmt"
	"log"

	// containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func listEvents(cli *client.Client, ctx context.Context) {
	eventFilters := filters.NewArgs()
	eventFilters.Add("type", "network")
	eventFilters.Add("event", "connect")
	eventFilters.Add("event", "disconnect")

	msgChan, errChan := cli.Events(ctx, events.ListOptions{Filters: eventFilters})
	for {
		select {
		case msg := <-msgChan:
			fmt.Println(msg.Actor)
		case err := <-errChan:
			fmt.Println(err)
		}
	}
}

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panic("Cannot initialize docker client")
	}
	defer cli.Close()

	listEvents(cli, ctx)
}
