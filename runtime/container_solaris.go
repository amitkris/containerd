package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/docker/containerd/specs"
	"github.com/opencontainers/runc/libcontainer"
	//	_ocs _"github.com/opencontainers/specs/specs-go"
)

func getRootIDs(s *specs.Spec) (int, int, error) {
	return 0, 0, nil
}

func (c *container) State() State {
	proc := c.processes["init"]
	if proc == nil {
		return Stopped
	}
	return proc.State()
}

func (c *container) Runtime() string {
	return c.runtime
}

func (c *container) Pause() error {
	return fmt.Errorf("Pause is not supported on Solaris\n")
}

func (c *container) Resume() error {
	return fmt.Errorf("Resume is not supported on Solaris\n")
}

func (c *container) Checkpoints() ([]Checkpoint, error) {
	return nil, fmt.Errorf("Checkpoints not supported on Solaris\n")
}

func (c *container) Checkpoint(cpt Checkpoint) error {
	return fmt.Errorf("Checkpoint is not supported on Solaris\n")
}

func (c *container) DeleteCheckpoint(name string) error {
	return fmt.Errorf("DeleteCheckpoint is not supported on Solaris\n")
}

func (c *container) Start(checkpoint string, s Stdio) (Process, error) {
	processRoot := filepath.Join(c.root, c.id, InitProcessID)
	if err := os.Mkdir(processRoot, 0755); err != nil {
		return nil, err
	}
	cmd := exec.Command("containerd-shim",
		c.id, c.bundle, c.runtime,
	)
	cmd.Dir = processRoot
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	spec, err := c.readSpec()
	if err != nil {
		return nil, err
	}
	config := &processConfig{
		checkpoint:  checkpoint,
		root:        processRoot,
		id:          InitProcessID,
		c:           c,
		stdio:       s,
		spec:        spec,
		processSpec: specs.ProcessSpec(spec.Process),
	}
	p, err := newProcess(config)
	if err != nil {
		return nil, err
	}
	if err := c.startCmd(InitProcessID, cmd, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *container) Exec(pid string, pspec specs.ProcessSpec, s Stdio) (pp Process, err error) {
	processRoot := filepath.Join(c.root, c.id, pid)
	if err := os.Mkdir(processRoot, 0755); err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			c.RemoveProcess(pid)
		}
	}()
	cmd := exec.Command("containerd-shim",
		c.id, c.bundle, c.runtime,
	)
	cmd.Dir = processRoot
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	spec, err := c.readSpec()
	if err != nil {
		return nil, err
	}
	config := &processConfig{
		exec:        true,
		id:          pid,
		root:        processRoot,
		c:           c,
		processSpec: pspec,
		spec:        spec,
		stdio:       s,
	}
	p, err := newProcess(config)
	if err != nil {
		return nil, err
	}
	if err := c.startCmd(pid, cmd, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *container) startCmd(pid string, cmd *exec.Cmd, p *process) error {
	if err := cmd.Start(); err != nil {
		if exErr, ok := err.(*exec.Error); ok {
			if exErr.Err == exec.ErrNotFound || exErr.Err == os.ErrNotExist {
				return fmt.Errorf("containerd-shim not installed on system")
			}
		}
		return err
	}
	if err := waitForStart(p, cmd); err != nil {
		return err
	}
	c.processes[pid] = p
	return nil
}

func (c *container) getLibctContainer() (libcontainer.Container, error) {
	return nil, nil
}

func (c *container) Pids() ([]int, error) {
	return nil, fmt.Errorf("Pids not implemented on Solaris\n")
}

func (c *container) Stats() (*Stat, error) {
	return nil, fmt.Errorf("Stats not implemented on Solaris\n")
}

func (c *container) OOM() (OOM, error) {
	return nil, nil
}

func waitForStart(p *process, cmd *exec.Cmd) error {
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

// isAlive checks if the shim that launched the container is still alive
func isAlive(cmd *exec.Cmd) (bool, error) {
	if err := syscall.Kill(cmd.Process.Pid, 0); err != nil {
		if err == syscall.ESRCH {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type oom struct {
	id      string
	root    string
	control *os.File
	eventfd int
}

func (o *oom) ContainerID() string {
	return o.id
}

func (o *oom) FD() int {
	return o.eventfd
}

func (o *oom) Flush() {
	buf := make([]byte, 8)
	syscall.Read(o.eventfd, buf)
}

func (o *oom) Removed() bool {
	_, err := os.Lstat(filepath.Join(o.root, "cgroup.event_control"))
	return os.IsNotExist(err)
}

func (o *oom) Close() error {
	err := syscall.Close(o.eventfd)
	if cerr := o.control.Close(); err == nil {
		err = cerr
	}
	return err
}

type message struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func readLogMessages(path string) ([]message, error) {
	var out []message
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for {
		var m message
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				return out, nil
			}
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}
