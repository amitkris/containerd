// +build solaris,cgo

package supervisor

/*
#include <port.h>
#include <poll.h>
#include <stdio.h>
#include <unistd.h>

int portAssociate(int port, int fd) {
	if (port_associate(port, PORT_SOURCE_FD, fd, POLLIN | POLLHUP, NULL) < 0) {
		return 1;
	}
}

int getFd(uintptr_t x) {
	return *(int *)x;
}
*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/Sirupsen/logrus"
	"github.com/docker/containerd/runtime"
)

//XXX Solaris

// NewMonitor 1
func NewMonitor() (*Monitor, error) {
	m := &Monitor{
		receivers: make(map[int]interface{}),
		exits:     make(chan runtime.Process, 1024),
		ooms:      make(chan string, 1024),
	}
	fd, err := C.port_create()
	if err != nil {
		return nil, err
	}

	m.epollFd = int(fd)
	go m.start()
	return m, nil
}

// Monitor 1
type Monitor struct {
	m         sync.Mutex
	receivers map[int]interface{}
	exits     chan runtime.Process
	ooms      chan string
	epollFd   int
}

// Exits 1
func (m *Monitor) Exits() chan runtime.Process {
	return m.exits
}

// OOMs 1
func (m *Monitor) OOMs() chan string {
	return m.ooms
}

// Monitor 1
func (m *Monitor) Monitor(p runtime.Process) error {
	m.m.Lock()
	defer m.m.Unlock()
	fd := p.ExitFD()
	if _, err := C.port_associate(C.int(m.epollFd), C.PORT_SOURCE_FD, C.uintptr_t(fd), C.POLLIN|C.POLLHUP, unsafe.Pointer(&fd)); err != nil {
		return err
	}
	EpollFdCounter.Inc(1)
	m.receivers[fd] = p
	return nil
}

// MonitorOOM 1
func (m *Monitor) MonitorOOM(c runtime.Container) error {
	return nil
	m.m.Lock()
	defer m.m.Unlock()
	o, err := c.OOM()
	if err != nil {
		return err
	}
	fd := o.FD()
	if _, err := C.port_associate(C.int(m.epollFd), C.PORT_SOURCE_FD, C.uintptr_t(fd), C.POLLIN|C.POLLHUP, unsafe.Pointer(&fd)); err != nil {
		return err
	}
	EpollFdCounter.Inc(1)
	m.receivers[fd] = o
	return nil
}

// Close 1
func (m *Monitor) Close() error {
	_, err := C.close(C.int(m.epollFd))
	return err
}

func (m *Monitor) start() {
	var ev C.port_event_t
	for {
		if C.port_get(C.int(m.epollFd), &ev, nil) < 0 {
			logrus.Warnf("containerd: epoll wait")
		}
		fd := int(C.getFd(C.uintptr_t(uintptr((ev.portev_user)))))
		m.m.Lock()
		r := m.receivers[fd]
		switch t := r.(type) {
		case runtime.Process:
			if ev.portev_events == C.POLLHUP {
				delete(m.receivers, fd)
				if err := t.Close(); err != nil {
					logrus.Warnf("containerd: close process IO")
				}
				EpollFdCounter.Dec(1)
				m.exits <- t
			}
		case runtime.OOM:
			// always flush the event fd
			t.Flush()
			if t.Removed() {
				delete(m.receivers, fd)
				// epoll will remove the fd from its set after it has been closed
				t.Close()
				EpollFdCounter.Dec(1)
			} else {
				m.ooms <- t.ContainerID()
			}
		}
		m.m.Unlock()
	}
}
