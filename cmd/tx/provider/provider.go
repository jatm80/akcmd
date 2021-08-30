package provider

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/provider/config"
	"github.com/ovrclk/akash/x/provider/types"
	"github.com/ovrclk/akcmd/client"
	"github.com/pkg/errors"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "provider",
		Desc: "Provider transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{createCMD(), updateCMD()},
	}

	return cmd
}

func createCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "Create a provider",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)

			cmd.AddArg("config-file", "config-file", true)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			// TODO: enable reading files with non-local URIs
			cfg, err := config.ReadConfigPath(args[0])
			if err != nil {
				err = errors.Wrapf(err, "ReadConfigPath err: %q", args[0])
				return err
			}

			msg := &types.MsgCreateProvider{
				Owner:      clientCtx.GetFromAddress().String(),
				HostURI:    cfg.Host,
				Info:       cfg.Info,
				Attributes: cfg.GetAttributes(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func updateCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "update",
		Desc: "Update a provider",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)

			cmd.AddArg("config-file", "config-file", true)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			cfg, err := config.ReadConfigPath(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateProvider{
				Owner:      clientCtx.GetFromAddress().String(),
				HostURI:    cfg.Host,
				Info:       cfg.Info,
				Attributes: cfg.GetAttributes(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}
