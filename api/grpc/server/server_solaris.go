package server

import (
	"errors"

	"github.com/docker/containerd/api/grpc/types"
	"github.com/docker/containerd/specs"
	"github.com/docker/containerd/supervisor"
	"golang.org/x/net/context"
)

// noop on Solaris (Checkpoints not supported)
func createContainerConfigCheckpoint(e *supervisor.StartTask, c *types.CreateContainerRequest) {
}

// TODO Solaris - may be able to completely factor out
func (s *apiServer) CreateCheckpoint(ctx context.Context, r *types.CreateCheckpointRequest) (*types.CreateCheckpointResponse, error) {
	return nil, errors.New("CreateCheckpoint() not supported on Solaris")
}

// TODO Solaris - may be able to completely factor out
func (s *apiServer) DeleteCheckpoint(ctx context.Context, r *types.DeleteCheckpointRequest) (*types.DeleteCheckpointResponse, error) {
	return nil, errors.New("DeleteCheckpoint() not supported on Solaris")
}

// TODO Solaris - may be able to completely factor out
func (s *apiServer) ListCheckpoint(ctx context.Context, r *types.ListCheckpointRequest) (*types.ListCheckpointResponse, error) {
	return nil, errors.New("ListCheckpoint() not supported on Solaris")
}

func (s *apiServer) Stats(ctx context.Context, r *types.StatsRequest) (*types.StatsResponse, error) {
	return nil, errors.New("Stats() not supported on Solaris")
}

func setUserFieldsInProcess(p *types.Process, oldProc specs.ProcessSpec) {
}

func setPlatformRuntimeProcessSpecUserFields(r *types.AddProcessRequest, process *specs.ProcessSpec) {
}
