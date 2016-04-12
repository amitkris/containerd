package supervisor

import "fmt"

func (s *Supervisor) handleTask(i Task) {
	fmt.Printf("in handle task\n")
	var err error
	switch t := i.(type) {
	case *AddProcessTask:
		fmt.Printf("a\n")
		err = s.addProcess(t)
	case *StartTask:
		fmt.Printf("b\n")
		err = s.start(t)
	case *DeleteTask:
		fmt.Printf("c\n")
		err = s.delete(t)
	case *ExitTask:
		fmt.Printf("d\n")
		err = s.exit(t)
	case *ExecExitTask:
		fmt.Printf("e\n")
		err = s.execExit(t)
	case *GetContainersTask:
		fmt.Printf("f\n")
		err = s.getContainers(t)
	case *SignalTask:
		fmt.Printf("g\n")
		err = s.signal(t)
	case *StatsTask:
		fmt.Printf("h\n")
		err = s.stats(t)
	case *UpdateTask:
		fmt.Printf("i\n")
		err = s.updateContainer(t)
	case *UpdateProcessTask:
		fmt.Printf("j\n")
		err = s.updateProcess(t)
	default:
		fmt.Printf("k\n")
		err = ErrUnknownTask
	}
	if err != errDeferedResponse {
		fmt.Printf("error is: %+v\n", err)
		i.ErrorCh() <- err
		close(i.ErrorCh())
	}
}
