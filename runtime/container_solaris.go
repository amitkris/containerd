package runtime

import (
	"github.com/docker/containerd/specs"
)

func (c *container) OOM() (OOM, error) {
	return nil, nil
}

func (c *container) UpdateResources(r *Resource) error {
	return nil
}

func getRootIDs(s *specs.Spec) (int, int, error) {
	return 0, 0, nil
}
