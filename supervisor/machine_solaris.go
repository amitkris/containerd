package supervisor

import (
	_ "errors"
)

type Machine struct {
	Cpus   int
	Memory int64
}

func CollectMachineInformation() (Machine, error) {
	m := Machine{}
	//return m, errors.New("supervisor CollectMachineInformation not implemented on Solaris")
	return m, nil
}
