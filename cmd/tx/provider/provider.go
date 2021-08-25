package provider

import (
	"github.com/gookit/gcli/v3"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "provider",
		Desc: "Provider transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	cmd.Add(createCMD())
	cmd.Add(updateCMD())

	return cmd
}

func createCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "Create a provider",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func updateCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "update",
		Desc: "Update a provider",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
