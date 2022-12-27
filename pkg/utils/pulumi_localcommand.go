package utils

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (u *Utils) PulumiExecuteCommand(ctx *pulumi.Context, commandName, command string) (*local.Command, error) {
	commandArgs := &local.CommandArgs{
		ArchivePaths: nil,
		AssetPaths:   nil,
		Create:       pulumi.String(command),
		Delete:       nil,
		Dir:          nil,
		Environment:  nil,
		Interpreter:  nil,
		Stdin:        nil,
		Triggers:     nil,
		Update:       nil,
	}

	return local.NewCommand(ctx, commandName, commandArgs)
}
