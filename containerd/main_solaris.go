package main

import (
	"time"

	"github.com/codegangsta/cli"
)

const (
	defaultStateDir     = "/run/containerd"
	defaultListenType   = "unix"
	defaultGRPCEndpoint = "/run/containerd/containerd.sock"
)

func appendPlatformFlags() {
	daemonFlags = append(daemonFlags, cli.StringFlag{
		Name:  "graphite-address",
		Usage: "Address of graphite server",
	})
}

func setAppBefore(app *cli.App) {
}

func reapProcesses() {
}

func processMetrics() {
}

func debugMetrics(interval time.Duration, graphiteAddr string) error {
	return nil
}
