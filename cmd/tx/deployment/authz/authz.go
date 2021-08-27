package authz

import "github.com/gookit/gcli/v3"

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "authz",
		Desc: "Deployment authorization transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{grantCMD(), revokeCMD()},
	}

	return cmd
}

func grantCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "grant",
		Desc: "Grant deposit deployment authorization to an address",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func revokeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "revoke",
		Desc: "Revoke deposit deployment authorization given to an address",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
