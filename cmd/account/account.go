package account

import (
	"github.com/gookit/gcli/v3"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "account",
		// allow color tag and {$cmd} will be replace to 'demo'
		Desc: "Transactions subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
