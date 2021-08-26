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
		Subs: []*gcli.Command{create.Cmd(), revokeCMD()},
	}

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
