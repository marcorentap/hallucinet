package mapper

import (
	"fmt"
	"github.com/marcorentap/hallucinet/backend/core"
)

type ContainerNameMapper struct {
	hctx    core.HallucinetContext
	mapping map[string]string
}

func NewContainerNameMapper(hctx core.HallucinetContext) *ContainerNameMapper {
	m := ContainerNameMapper{hctx: hctx, mapping: map[string]string{}}
	return &m
}

func (m *ContainerNameMapper) AddContainer(containerID string) {
	client := m.hctx.Client
	containerName := client.GetContainerName(containerID)
	m.mapping[containerID] = containerName + ".test"
}

func (m *ContainerNameMapper) RemoveContainer(containerID string) {
	delete(m.mapping, containerID)
}

func (m *ContainerNameMapper) ToHosts() string {
	return fmt.Sprintf("%v", m.mapping)
}
