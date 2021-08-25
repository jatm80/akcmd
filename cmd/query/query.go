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

func Cmd(app *gcli.App) *gcli.Command {
	cmd := &gcli.Command{
		Name: "query",
		// allow color tag and {$cmd} will be replace to 'demo'
		Desc: "Querying subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			app.Run([]string{"query", "-h"})
			return nil
		},
	}

	cmd.Add(audit.Cmd())
	cmd.Add(cert.Cmd())
	cmd.Add(deployment.Cmd())
	cmd.Add(escrow.Cmd())
	cmd.Add(market.Cmd())
	cmd.Add(provider.Cmd())

	return cmd
}
