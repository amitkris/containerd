package runtime

import (
	"github.com/docker/containerd/specs"
	"github.com/opencontainers/runc/libcontainer"
	//	_ocs _"github.com/opencontainers/specs/specs-go"
)

func getRootIDs(s *specs.Spec) (int, int, error) {
	return 0, 0, nil
}

func (c *container) getLibctContainer() (libcontainer.Container, error) {
	return nil, nil
}

func (c *container) OOM() (OOM, error) {
	return nil, nil
}

func (c *container) UpdateResources(r *Resource) error {
	return nil
}
