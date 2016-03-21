// +build solaris

package libcontainer

import (
	"github.com/opencontainers/runc/libcontainer/zones"
)

// State represents a running container's state
type State struct {
	BaseState

	// Platform specific fields below here
}

// A libcontainer container object.
//
// Each container is thread-safe within the same process. Since a container can
// be destroyed by a separate process, any function may return that the container
// was not found.
type Container interface {
	BaseContainer

	// Methods below here are platform specific

}

// XXX: Should collect the networking stats too.
// 	For the stubs, need config with slice of NICs.
func GetStats(id string) (*Stats, error) {
	stats := &Stats{}
	zstats := zones.Stats{}
	cpuUsage := zones.CpuUsage{
		TotalUsage: 5,
		/* XXX: currently only TotalUsage is consumed.
		PercpuUsage:		[]uint64 { 1, 2, 3 },
		UsageInKernelmode:	3,
		UsageInUsermode:	2,
		*/
	}
	cpuStats := zones.CpuStats{
		CpuUsage: cpuUsage,
		/* XXX: currently only CpuUsage.TotalUsage is consumed.
		ThrottlingData:	ThrottlingData {
			Periods:		100,
			ThrottledPeriods:	5,
			ThrottledTime:		1,
		},
		*/
	}
	memoryStats := zones.MemoryStats{
		Cache: 65536,
		Usage: zones.MemoryData{
			Usage:    32000000,
			MaxUsage: 64000000,
			Failcnt:  10,
		},
		/* XXX: currently only MemboriStats.Usage is consumed
		SwapUsage: MemoryData {
			Usage:		8192000,
			MaxUsage:	8192000,
			Failcnt:	128,
		},
		KernelUsage: MemoryData {
			Usage:		4096000,
			MaxUsage:	2048000,
			Failcnt:	0,
		},
		*/
	}
	blkioStats := zones.BlkioStats{
		IoServiceBytesRecursive: []zones.BlkioStatEntry{
			{
				Major: 14,
				Minor: 1,
				Op:    "read", //op name from api/client/stats.go
				Value: 9000000,
			},
			{
				Major: 13,
				Minor: 0,
				Op:    "write", //op name from api/client/stats.go
				Value: 500000,
			},
		},
		/* XXX: currently only IoServiceBytesRecursive is consumed
		IoServicedRecursive: []BlkioStatEntry {
			{
				Major:	14000000,
				Minor:	10000000,
				Op:	"",
				Value:	9000000,
			},
		},
		*/
	}
	zstats.CpuStats = cpuStats
	zstats.MemoryStats = memoryStats
	zstats.BlkioStats = blkioStats

	stats.Stats = &zstats
	return stats, nil
}
