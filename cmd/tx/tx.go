package tx

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/tx/audit"
	"github.com/ovrclk/akcmd/cmd/tx/cert"
	"github.com/ovrclk/akcmd/cmd/tx/deployment"
	"github.com/ovrclk/akcmd/cmd/tx/market"
	"github.com/ovrclk/akcmd/cmd/tx/provider"
)

func Cmd(app *gcli.App) *gcli.Command {
	cmd := &gcli.Command{
		Name: "tx",
		// allow color tag and {$cmd} will be replace to 'demo'
		Desc: "Transactions subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			app.Run([]string{"tx", "-h"})
			return nil
		},
	}

	cmd.Add(audit.Cmd())
	cmd.Add(cert.Cmd())
	cmd.Add(deployment.Cmd())
	cmd.Add(market.Cmd())
	cmd.Add(provider.Cmd())

	return cmd
}
