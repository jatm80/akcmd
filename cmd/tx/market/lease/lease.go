package lease

import "github.com/gookit/gcli/v3"

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "lease",
		Desc: "Lease subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{createCMD(), withdrawCMD(), closeCMD()},
	}

	return cmd
}

func createCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "Create a market lease",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func withdrawCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "withdraw",
		Desc: "Settle and withdraw available funds from market order escrow account",
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
		Desc: "Close a market lease",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
