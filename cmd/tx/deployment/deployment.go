package deployment

import (
	"context"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/sdl"
	cutils "github.com/ovrclk/akash/x/cert/utils"
	"github.com/ovrclk/akash/x/deployment/client/cli"
	"github.com/ovrclk/akash/x/deployment/types"
	"github.com/ovrclk/akcmd/client"
	"github.com/ovrclk/akcmd/cmd/tx/deployment/authz"
	"github.com/ovrclk/akcmd/cmd/tx/deployment/group"
	"github.com/pkg/errors"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "deployment",
		Desc: "Deployment transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{
			createCMD(), updateCMD(), depositCMD(), closeCMD(),
			group.Cmd(), authz.Cmd(),
		},
	}

	return cmd
}

func createCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "Create deployment",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddDeploymentIDFlags(cmd)
			client.AddDepositFlags(cmd, cli.DefaultDeposit)
			client.AddDepositorFlag(cmd)

			cmd.AddArg("sdl-file", "sdl-file", true)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			// first lets validate certificate exists for given account
			if _, err = cutils.LoadAndQueryCertificateForAccount(context.Background(), clientCtx, clientCtx.Keyring); err != nil {
				if os.IsNotExist(err) {
					err = errors.Errorf("no certificate file found for account %q.\n"+
						"consider creating it as certificate required to submit manifest", clientCtx.FromAddress.String())
				}

				return err
			}

			sdlManifest, err := sdl.ReadFile(args[0])
			if err != nil {
				return err
			}

			groups, err := sdlManifest.DeploymentGroups()
			if err != nil {
				return err
			}

			id, err := client.DeploymentIDFromFlags(cli.WithOwner(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			// Default DSeq to the current block height
			if id.DSeq == 0 {
				if id.DSeq, err = cli.CurrentBlockHeight(clientCtx); err != nil {
					return err
				}
			}

			version, err := sdl.Version(sdlManifest)
			if err != nil {
				return err
			}

			deposit, err := client.DepositFromFlags()
			if err != nil {
				return err
			}

			depositorAcc, err := client.DepositorFromFlags(id.Owner)
			if err != nil {
				return err
			}

			msg := &types.MsgCreateDeployment{
				ID:        id,
				Version:   version,
				Groups:    make([]types.GroupSpec, 0, len(groups)),
				Deposit:   deposit,
				Depositor: depositorAcc,
			}

			for _, group := range groups {
				msg.Groups = append(msg.Groups, *group)
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
		Desc: "Update deployment",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddDeploymentIDFlags(cmd)

			cmd.AddArg("sdl-file", "sdl-file", true)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.DeploymentIDFromFlags(cli.WithOwner(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			sdlManifest, err := sdl.ReadFile(args[0])
			if err != nil {
				return err
			}
			groups, err := sdlManifest.DeploymentGroups()
			if err != nil {
				return err
			}

			version, err := sdl.Version(sdlManifest)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateDeployment{
				ID:      id,
				Version: version,
				Groups:  make([]types.GroupSpec, 0, len(groups)),
			}

			for _, group := range groups {
				msg.Groups = append(msg.Groups, *group)
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func depositCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "deposit",
		Desc: "Deposit deployment",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddDeploymentIDFlags(cmd)
			client.AddDepositorFlag(cmd)

			cmd.AddArg("amount", "deposit amount", true)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.DeploymentIDFromFlags(cli.WithOwner(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			depositorAcc, err := client.DepositorFromFlags(id.Owner)
			if err != nil {
				return err
			}

			msg := &types.MsgDepositDeployment{
				ID:        id,
				Amount:    deposit,
				Depositor: depositorAcc,
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func closeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "close",
		Desc: "Close deployment",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddDeploymentIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.DeploymentIDFromFlags(cli.WithOwner(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			msg := &types.MsgCloseDeployment{ID: id}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}
