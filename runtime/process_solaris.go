package runtime

import (
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

// Signal sends the provided signal to the process
func (p *process) Signal(s os.Signal) error {
	return syscall.Kill(p.pid, s.(syscall.Signal))
}

func populateProcessStateForEncoding(config *processConfig, uid int, gid int) ProcessState {
	return ProcessState{
		PlatformProcessState: PlatformProcessState{
			Checkpoint: config.checkpoint,
			RootUID:    uid,
			RootGID:    gid,
		},
		Stdin:  config.stdio.Stdin,
		Stdout: config.stdio.Stdout,
		Stderr: config.stdio.Stderr,
	}
}

func getExitPipe(path string) (*os.File, error) {
	if err := unix.Mkfifo(path, 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}
	// add NONBLOCK in case the other side has already closed or else
	// this function would never return
	return os.OpenFile(path, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
}

func getControlPipe(path string) (*os.File, error) {
	if err := unix.Mkfifo(path, 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}
	return os.OpenFile(path, syscall.O_RDWR|syscall.O_NONBLOCK, 0)
}
