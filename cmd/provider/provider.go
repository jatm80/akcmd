package provider

import (
	"github.com/gookit/gcli/v3"
)

func Cmd(app *gcli.App) *gcli.Command {
	cmd := &gcli.Command{
		Name: "provider",
		// allow color tag and {$cmd} will be replace to 'demo'
		Desc: "Query and interact with providers on the Network",
		Func: func(cmd *gcli.Command, args []string) error {
			app.Run([]string{"provider", "-h"})
			return nil
		},
	}

	cmd.Add(sendManifestCMD())
	cmd.Add(statusCMD())
	cmd.Add(leaseStatusCMD())
	cmd.Add(leaseEventsCMD())
	cmd.Add(leaseLogsCMD())
	cmd.Add(serviceStatusCMD())
	cmd.Add(leaseShellCMD())

	return cmd
}

func sendManifestCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "send-manifest",
		Desc: "Submit manifest to provider(s)",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func statusCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "status",
		Desc: "Get provider status",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func leaseStatusCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "lease-status",
		Desc: "Get lease status",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func leaseEventsCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "lease-events",
		Desc: "Get lease events",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func leaseLogsCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "lease-logs",
		Desc: "Get lease logs",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func serviceStatusCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "service-status",
		Desc: "Get service status",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}

func leaseShellCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "lease-shell",
		Desc: "Do lease shell",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
