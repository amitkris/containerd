package supervisor

type Machine struct {
	Cpus   int
	Memory int64
}

func CollectMachineInformation() (Machine, error) {
	m := Machine{}
	return m, nil
}
