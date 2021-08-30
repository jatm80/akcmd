package authz

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/deployment/types"
	"github.com/ovrclk/akcmd/client"
	"github.com/pkg/errors"
)

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
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddExpirationFlag(cmd)
			cmd.Required(flags.FlagFrom)

			cmd.AddArg("grantee", "grantee", true)
			cmd.AddArg("spend_limit", "spend_limit", true)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			grantee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			spendLimit, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			if spendLimit.IsZero() || spendLimit.IsNegative() {
				return errors.Errorf("spend-limit should be greater than zero, got: %s", spendLimit)
			}

			exp := client.ExpirationFromFlag()

			granter := clientCtx.GetFromAddress()
			authorization := types.NewDepositDeploymentAuthorization(spendLimit)

			msg, err := authz.NewMsgGrant(granter, grantee, authorization, time.Unix(exp, 0))
			if err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func revokeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "revoke",
		Desc: "Revoke deposit deployment authorization given to an address",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			cmd.Required(flags.FlagFrom)

			cmd.AddArg("grantee", "grantee", true)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			grantee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			granter := clientCtx.GetFromAddress()
			msgTypeURL := types.DepositDeploymentAuthorization{}.MsgTypeURL()
			msg := authz.NewMsgRevoke(granter, grantee, msgTypeURL)

			return client.BroadcastTX(clientCtx, &msg)
		},
	}

	return cmd
}
