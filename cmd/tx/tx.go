package tx

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/tx/audit"
	"github.com/ovrclk/akcmd/cmd/tx/cert"
	"github.com/ovrclk/akcmd/cmd/tx/deployment"
	"github.com/ovrclk/akcmd/cmd/tx/market"
	"github.com/ovrclk/akcmd/cmd/tx/provider"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "tx",
		// allow color tag and {$cmd} will be replace to 'demo'
		Desc: "Transactions subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{
			audit.Cmd(), cert.Cmd(), deployment.Cmd(), market.Cmd(), provider.Cmd(),
		},
	}

	return cmd
}
