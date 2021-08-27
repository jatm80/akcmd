package config

import (
	"bytes"
	"log"
	"os"

	"github.com/gookit/color"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/interact"

	"github.com/ovrclk/akcmd/l10n"
)

type Options struct {
	ConfigFileName string
	ConfigFilePath string
}

func getHome() string {
	dirname, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	return dirname
}

func writeConfig(path string, options Options) {
	configPath := path + options.ConfigFileName
	f, err := os.Create(configPath)

	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	_, err = config.DumpTo(buf, config.Yaml)

	if err != nil {
		panic(err)
	}

	color.Info.Printf("Saving user config file %s\n\n", configPath)
	f.WriteString(buf.String())
}

func createInitConfig() {
	config.SetData(map[string]interface{}{
		"name":    "akash_lite_config",
		"version": "0.12.2",
	})
}

func createConfigPath(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

func loadConfig(path string, options Options) bool {
	configPath := path + options.ConfigFileName
	err := config.LoadFiles(configPath)
	return err == nil
}

// returns if the application was first run
func LoadConfig(app *gcli.App, options Options) bool {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	homeDir := getHome()
	clientDir := homeDir + options.ConfigFilePath
	configLoaded := loadConfig(clientDir, options)

	if !configLoaded {
		app.Run([]string{"welcome"})
		createConfigPath(clientDir)
		createInitConfig()

		// walkthrough some stuff for user
		app.Run([]string{"config", "network", "update"})

		writeConfig(clientDir, options)
		return true
	}

	return false
}

func Cmd(options Options) *gcli.Command {
	localizedStrings := l10n.GetLocalizationStrings()

	cmd := &gcli.Command{
		Name: "config",
		Desc: localizedStrings.Command["config"],
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{
			{
				Name: "network",
				Desc: localizedStrings.Command["config.network"],
				// Func: func (cmd *gcli.Command, args []string) error {
				// 	GetNetworkInfo()
				// 	return nil
				// },
				Subs: []*gcli.Command{
					{
						Name: "update",
						Desc: localizedStrings.Command["config.network.update"],
						Func: func(cmd *gcli.Command, args []string) error {
							networkNames := []string{}

							// populate default network selection
							for i := range NetworkRegisteries {
								registry := NetworkRegisteries[i]
								networkNames = append(networkNames, registry.Name)
							}

							// select default network
							_ = interact.SelectOne(
								localizedStrings.Command["config.network.defaultSelection"],
								networkNames,
								"",
							)

							GetNetworkInfo()

							return nil
						},
					},
				},
			},
		},
	}

	return cmd
}
