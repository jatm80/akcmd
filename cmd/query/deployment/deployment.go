package deployment

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/query/deployment/group"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "deployment",
		Desc: "Deployment query commands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	cmd.Add(listCMD())
	cmd.Add(getCMD())
	cmd.Add(group.Cmd())

	return cmd
}

func listCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "list",
		Desc: "Query for all deployments",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func getCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "get",
		Desc: "Query deployment",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
