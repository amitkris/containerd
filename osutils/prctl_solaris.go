// +build solaris

package osutils

//Solaris TODO
// GetSubreaper returns the subreaper setting for the calling process
func GetSubreaper() (int, error) {
	return 0, nil
}

// SetSubreaper sets the value i as the subreaper setting for the calling process
func SetSubreaper(i int) error {
	return nil
}
