package mapper

import (
	"fmt"
	"strings"

	"github.com/marcorentap/hallucinet/backend/core"
)

type mappingEntry struct {
	domainName string
	address    string
}

type ContainerNameMapper struct {
	hctx    core.HallucinetContext
	mapping map[string]mappingEntry
}

func NewContainerNameMapper(hctx core.HallucinetContext) *ContainerNameMapper {
	m := ContainerNameMapper{hctx: hctx, mapping: map[string]mappingEntry{}}
	return &m
}

func (m *ContainerNameMapper) AddContainer(containerID string) {
	client := m.hctx.Client
	containerName := client.GetContainerName(containerID)
	containerAddr := client.GetContainerAddr(containerID)
	m.mapping[containerID] = mappingEntry{
		domainName: containerName[1:] + ".test",
		address:    containerAddr,
	}
}

func (m *ContainerNameMapper) RemoveContainer(containerID string) {
	delete(m.mapping, containerID)
}

func (m *ContainerNameMapper) ToHosts() string {
	var s strings.Builder
	for _, entry := range m.mapping {
		s.WriteString(fmt.Sprintf("%v %v\n", entry.address, entry.domainName))
	}
	return s.String()
}
