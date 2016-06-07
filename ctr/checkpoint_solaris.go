package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var checkpointSubCmds = []cli.Command{
	listCheckpointCommand,
	createCheckpointCommand,
	deleteCheckpointCommand,
}

var checkpointCommand = cli.Command{
	Name:        "checkpoints",
	Usage:       "list all checkpoints",
	ArgsUsage:   "COMMAND [arguments...]",
	Subcommands: checkpointSubCmds,
	Description: func() string {
		desc := "\n    COMMAND:\n"
		for _, command := range checkpointSubCmds {
			desc += fmt.Sprintf("    %-10.10s%s\n", command.Name, command.Usage)
		}
		return desc
	}(),
	Action: listCheckpoints,
}

var listCheckpointCommand = cli.Command{
	Name:   "list",
	Usage:  "list all checkpoints for a container",
	Action: listCheckpoints,
}

func listCheckpoints(context *cli.Context) {
	fatal("checkpoint command is not supported on Solaris", ExitStatusUnsupported)
}

var createCheckpointCommand = cli.Command{
	Name:  "create",
	Usage: "create a new checkpoint for the container",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "tcp",
			Usage: "persist open tcp connections",
		},
		cli.BoolFlag{
			Name:  "unix-sockets",
			Usage: "perist unix sockets",
		},
		cli.BoolFlag{
			Name:  "exit",
			Usage: "exit the container after the checkpoint completes successfully",
		},
		cli.BoolFlag{
			Name:  "shell",
			Usage: "checkpoint shell jobs",
		},
	},
	Action: func(context *cli.Context) {
		fatal("checkpoint command is not supported on Solaris", ExitStatusUnsupported)
	},
}

var deleteCheckpointCommand = cli.Command{
	Name:  "delete",
	Usage: "delete a container's checkpoint",
	Action: func(context *cli.Context) {
		fatal("checkpoint command is not supported on Solaris", ExitStatusUnsupported)
	},
}
