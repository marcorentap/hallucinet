package types

import (
	"database/sql"
	"time"
)

type HallucinetConfig struct {
	NetworkName  string
	SqlitePath   string
	DomainSuffix string
	HostsPath    string
}

type HallucinetContext struct {
	Config    HallucinetConfig
	EventChan chan HallucinetEvent
	DB        *sql.DB
}

type HallucinetEvent struct {
	Kind          HallucinetEventKind
	ContainerIP   string
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
