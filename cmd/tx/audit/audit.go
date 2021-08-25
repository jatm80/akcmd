package audit

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/tx/audit/attr"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "audit",
		Desc: "Audit transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	cmd.Add(attr.Cmd())

	return cmd
}
