package create

import "github.com/gookit/gcli/v3"

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "create/update api certificates",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	cmd.Add(clientCMD())
	cmd.Add(serverCMD())

	return cmd
}

func clientCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "client",
		Desc: "create client api certificate",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func serverCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "server",
		Desc: "create server api certificate",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
