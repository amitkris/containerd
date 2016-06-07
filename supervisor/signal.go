package supervisor

import (
	"fmt"
	"os"
)

type SignalTask struct {
	baseTask
	ID     string
	PID    string
	Signal os.Signal
}

func (s *Supervisor) signal(t *SignalTask) error {
	fmt.Printf("in supervisor signal\n")
	i, ok := s.containers[t.ID[0:12]]
	if !ok {
		return ErrContainerNotFound
	}
	processes, err := i.container.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		fmt.Printf("sending signal to pid: %+v\n", p.ID())
		if p.ID() == t.PID {
			return p.Signal(t.Signal)
		}
	}
	return ErrProcessNotFound
}
