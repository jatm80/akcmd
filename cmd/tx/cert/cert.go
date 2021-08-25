package cert

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/tx/cert/create"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "cert",
		Desc: "Certificates transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	cmd.Add(create.Cmd())
	cmd.Add(revokeCMD())

	return cmd
}

func revokeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "revoke",
		Desc: "revoke api certificate",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
