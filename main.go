package main

import (
	"github.com/gookit/gcli/v3"

	account "github.com/ovrclk/akcmd/cmd/account"
	config "github.com/ovrclk/akcmd/cmd/config"
	deploy "github.com/ovrclk/akcmd/cmd/deploy"
	node "github.com/ovrclk/akcmd/cmd/node"
	project "github.com/ovrclk/akcmd/cmd/project"
	provider "github.com/ovrclk/akcmd/cmd/provider"
	query "github.com/ovrclk/akcmd/cmd/query"
	tx "github.com/ovrclk/akcmd/cmd/tx"
	welcome "github.com/ovrclk/akcmd/cmd/welcome"

	"github.com/ovrclk/akcmd/l10n"
)

func main() {
	localizedStrings := l10n.GetLocalizationStrings()

	app := gcli.NewApp()
	app.Version = "0.12.2"
	app.Desc = localizedStrings.AppDescription
	app.ExitOnEnd = false
	// app.SetVerbose(gcli.VerbDebug)
	// app.Add(builtin.GenAutoComplete()

	// options
	configOptions := config.Options{
		ConfigFileName: `/akcmd_config.yml`,
		ConfigFilePath: `/.akash`,
	}

	// add commands to the CLI
	app.Add(account.Cmd(app))
	app.Add(config.Cmd(app, configOptions))
	app.Add(deploy.Cmd(app))
	app.Add(node.Cmd(app))
	app.Add(project.Cmd(app))
	app.Add(provider.Cmd(app))
	app.Add(query.Cmd(app))
	app.Add(tx.Cmd(app))
	app.Add(versionCMD(app))
	app.Add(welcome.Cmd(app))

	config.LoadConfig(app, configOptions)

	app.Run(nil)
}

func versionCMD(app *gcli.App) *gcli.Command {
	cmd := &gcli.Command{
		Name:    "version",
		Desc:    "Print the application binary version information",
		Aliases: []string{"v"},
		Func: func(cmd *gcli.Command, args []string) error {
			app.Run([]string{"--version"})
			return nil
		},
	}

	return cmd
}
