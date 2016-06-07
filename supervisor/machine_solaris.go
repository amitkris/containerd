package supervisor

/*
#include <unistd.h>
*/
import "C"

import (
	"errors"
)

type Machine struct {
	Cpus   int
	Memory int64
}

func CollectMachineInformation() (Machine, error) {
	m := Machine{}
	ncpus := C.sysconf(C._SC_NPROCESSORS_ONLN)
	if ncpus <= 0 {
		return m, errors.New("Unable to get number of cpus")
	}
	m.Cpus = int(ncpus)

	memTotal := getTotalMem()
	if memTotal < 0 {
		return m, errors.New("Unable to get total memory")
	}
	m.Memory = int64(memTotal / 1024 / 1024)
	return m, nil
}

// Get the system memory info using sysconf same as prtconf
func getTotalMem() int64 {
	pagesize := C.sysconf(C._SC_PAGESIZE)
	npages := C.sysconf(C._SC_PHYS_PAGES)
	return int64(pagesize * npages)
}
