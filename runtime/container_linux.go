package runtime

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/docker/containerd/specs"
	"github.com/opencontainers/runc/libcontainer"
)

func (c *container) getLibctContainer() (libcontainer.Container, error) {
	runtimeRoot := "/run/runc"

	// Check that the root wasn't changed
	for _, opt := range c.runtimeArgs {
		if strings.HasPrefix(opt, "--root=") {
			runtimeRoot = strings.TrimPrefix(opt, "--root=")
			break
		}
	}

	f, err := libcontainer.New(runtimeRoot, libcontainer.Cgroupfs)
	if err != nil {
		return nil, err
	}
	return f.Load(c.id)
}

func (c *container) OOM() (OOM, error) {
	container, err := c.getLibctContainer()
	if err != nil {
		if lerr, ok := err.(libcontainer.Error); ok {
			// with oom registration sometimes the container can run, exit, and be destroyed
			// faster than we can get the state back so we can just ignore this
			if lerr.Code() == libcontainer.ContainerNotExists {
				return nil, ErrContainerExited
			}
		}
		return nil, err
	}
	state, err := container.State()
	if err != nil {
		return nil, err
	}
	memoryPath := state.CgroupPaths["memory"]
	return c.getMemeoryEventFD(memoryPath)
}

func (c *container) UpdateResources(r *Resource) error {
	container, err := c.getLibctContainer()
	if err != nil {
		return err
	}
	config := container.Config()
	config.Cgroups.Resources.CpuShares = r.CPUShares
	config.Cgroups.Resources.BlkioWeight = r.BlkioWeight
	config.Cgroups.Resources.CpuPeriod = r.CPUPeriod
	config.Cgroups.Resources.CpuQuota = r.CPUQuota
	config.Cgroups.Resources.CpusetCpus = r.CpusetCpus
	config.Cgroups.Resources.CpusetMems = r.CpusetMems
	config.Cgroups.Resources.KernelMemory = r.KernelMemory
	config.Cgroups.Resources.Memory = r.Memory
	config.Cgroups.Resources.MemoryReservation = r.MemoryReservation
	config.Cgroups.Resources.MemorySwap = r.MemorySwap
	return container.Set(config)
}

func getRootIDs(s *specs.Spec) (int, int, error) {
	if s == nil {
		return 0, 0, nil
	}
	var hasUserns bool
	for _, ns := range s.Linux.Namespaces {
		if ns.Type == ocs.UserNamespace {
			hasUserns = true
			break
		}
	}
	if !hasUserns {
		return 0, 0, nil
	}
	uid := hostIDFromMap(0, s.Linux.UIDMappings)
	gid := hostIDFromMap(0, s.Linux.GIDMappings)
	return uid, gid, nil
}

func (c *container) getMemeoryEventFD(root string) (*oom, error) {
	f, err := os.Open(filepath.Join(root, "memory.oom_control"))
	if err != nil {
		return nil, err
	}
	fd, _, serr := syscall.RawSyscall(syscall.SYS_EVENTFD2, 0, syscall.FD_CLOEXEC, 0)
	if serr != 0 {
		f.Close()
		return nil, serr
	}
	if err := c.writeEventFD(root, int(f.Fd()), int(fd)); err != nil {
		syscall.Close(int(fd))
		f.Close()
		return nil, err
	}
	return &oom{
		root:    root,
		id:      c.id,
		eventfd: int(fd),
		control: f,
	}, nil
}
