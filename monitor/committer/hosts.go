package committer

import (
	"github.com/marcorentap/hallucinet/core"
	"log"
	"os"
)

type HostsCommitter struct {
	hostsPath string
	hctx      core.HallucinetContext
}

func NewHostsCommitter(hctx core.HallucinetContext) *HostsCommitter {
	return &HostsCommitter{
		hostsPath: hctx.Config.HostsPath,
		hctx:      hctx,
	}
}

func (c *HostsCommitter) Commit() {
	file, createErr := os.Create(c.hostsPath)
	if createErr != nil {
		log.Panicf("Cannot create or truncate file %v: %v\n", c.hostsPath, createErr)
	}

	_, writeErr := file.WriteString(c.hctx.Mapper.ToHosts())
	if writeErr != nil {
		log.Panicf("Cannot write to hosts file %v: %v\n", c.hostsPath, writeErr)
	}
}
