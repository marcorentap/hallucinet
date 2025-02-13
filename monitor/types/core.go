package types

import "time"

type HallucinetConfig struct {
	NetworkName string
}

type HallucinetContext struct {
	Config    HallucinetConfig
	EventChan chan HallucinetEvent
}

type HallucinetEvent struct {
	Kind          HallucinetEventKind
	ContainerID   string
	ContainerName string
	NetworkID     string
	NetworkName   string
	Data          any
	Time          time.Time
}

type HallucinetEventKind int

const (
	ContainerConnected = iota
	ContainerDisconnected
	UnknownEvent
)
