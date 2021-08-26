package query

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/query/audit"
	"github.com/ovrclk/akcmd/cmd/query/cert"
	"github.com/ovrclk/akcmd/cmd/query/deployment"
	"github.com/ovrclk/akcmd/cmd/query/escrow"
	"github.com/ovrclk/akcmd/cmd/query/market"
	"github.com/ovrclk/akcmd/cmd/query/provider"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "query",
		// allow color tag and {$cmd} will be replace to 'demo'
		Desc: "Querying subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{
			audit.Cmd(), cert.Cmd(), deployment.Cmd(), escrow.Cmd(),
			market.Cmd(), provider.Cmd(),
		},
	}

	return cmd
}
