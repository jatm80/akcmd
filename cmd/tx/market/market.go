package market

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/tx/market/bid"
	"github.com/ovrclk/akcmd/cmd/tx/market/lease"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "market",
		Desc: "Market transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	cmd.Add(bid.Cmd())
	cmd.Add(lease.Cmd())

	return cmd
}
