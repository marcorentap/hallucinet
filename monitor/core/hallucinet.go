package core

type ContainerClient interface {
	GetContainers() []string
	GetContainerName(containerID string) string
	GetContainerAddr(containerID string) string
	GetEvents() <-chan ContainerEvent
}

type Mapper interface {
	AddContainer(containerID string)
	RemoveContainer(containerID string)
	ToHosts() string
}

type Committer interface {
	Commit()
}

type ContainerEvent struct {
	ContainerID string
	EventType   ContainerEventType
}

type ContainerEventType int

const (
	EventConnected ContainerEventType = iota
	EventDisconnected
	EventUnknown
)

type HallucinetConfig struct {
	Client      string
	Mapper      string
	NetworkName string
	Committer   string
	Suffix      string
	HostsPath   string
}

type HallucinetContext struct {
	Config    HallucinetConfig
	Client    ContainerClient
	Mapper    Mapper
	Committer Committer
}
