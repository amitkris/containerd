// +build solaris

package main

import (
	"os"
)

// NewConsole returns an initalized console that can be used within a container by copying bytes
// from the master side to the slave that is attached as the tty for the container's init process.
func newConsole(uid, gid int) (*os.File, string, error) {
	return nil, "", nil
}
