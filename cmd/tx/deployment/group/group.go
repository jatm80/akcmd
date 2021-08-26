package group

import "github.com/gookit/gcli/v3"

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "group",
		Desc: "Modify a Deployment's specific Group",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{closeCMD(), pauseCMD(), startCMD()},
	}

	return cmd
}

func closeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "close",
		Desc: "close a Deployment's specific Group",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func pauseCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "pause",
		Desc: "pause a Deployment's specific Group",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func startCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "start",
		Desc: "start a Deployment's specific Group",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
