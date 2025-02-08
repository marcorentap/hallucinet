package mapper

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/client"
)

type ContainerNameMapper struct {
	ctx     context.Context
	cli     *client.Client
	mapping map[string]string
}

func NewContainerNameMapper(ctx context.Context, cli *client.Client) *ContainerNameMapper {
	return &ContainerNameMapper{
		ctx:     ctx,
		cli:     cli,
		mapping: map[string]string{},
	}
}

func (m *ContainerNameMapper) getContainerName(containerID string) string {
	con, err := m.cli.ContainerInspect(m.ctx, containerID)
	if err != nil {
		log.Panicf("Cannot inspect container %v: %v", containerID, err)
	}
	return con.Name
}

func (m *ContainerNameMapper) AddContainer(containerID string) {
	containerName := strings.TrimPrefix(m.getContainerName(containerID), "/")
	m.mapping[containerID] = containerName + ".test"
}

func (m *ContainerNameMapper) RemoveContainer(containerID string) {
	delete(m.mapping, containerID)
}

func (m *ContainerNameMapper) ToHosts() {
	log.Printf("%v", m.mapping)
}
