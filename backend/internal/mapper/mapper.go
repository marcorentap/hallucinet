package mapper

type Mapper interface {
	AddContainer(containerID string)
	RemoveContainer(containerID string)
	ToHosts()
}
