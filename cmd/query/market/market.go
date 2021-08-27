package market

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/query/market/bid"
	"github.com/ovrclk/akcmd/cmd/query/market/lease"
	"github.com/ovrclk/akcmd/cmd/query/market/order"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "market",
		Desc: "Market query commands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{order.Cmd(), bid.Cmd(), lease.Cmd()},
	}

	return cmd
}
