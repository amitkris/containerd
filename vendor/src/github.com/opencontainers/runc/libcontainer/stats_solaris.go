package libcontainer

import "github.com/opencontainers/runc/libcontainer/zones"

type Stats struct {
	Interfaces []*NetworkInterface
	Stats      *zones.Stats
}
