package runtime

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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

func waitForStart(p *process, cmd *exec.Cmd) error {
	return nil

	for i := 0; i < 300; i++ {
		if _, err := p.getPidFromFile(); err != nil {
			if os.IsNotExist(err) || err == errInvalidPidInt {
				alive, err := isAlive(cmd)
				if err != nil {
					return err
				}
				if !alive {
					// runc could have failed to run the container so lets get the error
					// out of the logs or the shim could have encountered an error
					messages, err := readLogMessages(filepath.Join(p.root, "shim-log.json"))
					if err != nil {
						return err
					}
					for _, m := range messages {
						if m.Level == "error" {
							return errors.New(m.Msg)
						}
					}
					// no errors reported back from shim, check for runc/runtime errors
					messages, err = readLogMessages(filepath.Join(p.root, "log.json"))
					if err != nil {
						if os.IsNotExist(err) {
							return ErrContainerNotStarted
						}
						return err
					}
					for _, m := range messages {
						if m.Level == "error" {
							return errors.New(m.Msg)
						}
					}
					return ErrContainerNotStarted
				}
				time.Sleep(50 * time.Millisecond)
				continue
			}
			return err
		}
		return nil
	}
	return errNoPidFile
}

func (c *container) UpdateResources(r *Resource) error {
	return nil
}
