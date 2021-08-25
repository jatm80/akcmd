package deployment

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/cmd/tx/deployment/authz"
	"github.com/ovrclk/akcmd/cmd/tx/deployment/group"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "deployment",
		Desc: "Deployment transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	cmd.Add(createCMD())
	cmd.Add(updateCMD())
	cmd.Add(depositCMD())
	cmd.Add(closeCMD())
	cmd.Add(group.Cmd())
	cmd.Add(authz.Cmd())

	return cmd
}

func createCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "Create deployment",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func updateCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "update",
		Desc: "Update deployment",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func depositCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "deposit",
		Desc: "Deposit deployment",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func closeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "close",
		Desc: "Close deployment",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
