package command

import (
	"github.com/gookit/gcli/v3"
)

func Cmd(app *gcli.App) *gcli.Command {
	cmd := &gcli.Command{
		Name: "deploy",
		// allow color tag and {$cmd} will be replace to 'demo'
		Desc: "Interact with deployments within the Akash Project",
		Func: func(cmd *gcli.Command, args []string) error {
			app.Run([]string{"deploy", "-h"})
			return nil
		},
	}

	cmd.Add(createCMD())

	return cmd
}

func createCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "Create a deployment on the akash network",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
