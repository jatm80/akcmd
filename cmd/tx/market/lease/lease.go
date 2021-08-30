package lease

import (
	"github.com/gookit/gcli/v3"
	dcli "github.com/ovrclk/akash/x/deployment/client/cli"
	"github.com/ovrclk/akash/x/market/types"
	"github.com/ovrclk/akcmd/client"
)

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
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddLeaseIDFlags(cmd)
			client.MarkReqLeaseIDFlags(cmd, dcli.DeploymentIDOptionNoOwner(true))
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.LeaseIDFromFlags(dcli.WithOwner(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			msg := &types.MsgCreateLease{
				BidID: id.BidID(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func withdrawCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "withdraw",
		Desc: "Settle and withdraw available funds from market order escrow account",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddLeaseIDFlags(cmd)
			client.MarkReqLeaseIDFlags(cmd, dcli.DeploymentIDOptionNoOwner(true))
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.LeaseIDFromFlags(dcli.WithOwner(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			msg := &types.MsgWithdrawLease{
				LeaseID: id,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func closeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "close",
		Desc: "Close a market lease",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddLeaseIDFlags(cmd)
			client.MarkReqLeaseIDFlags(cmd, dcli.DeploymentIDOptionNoOwner(true))
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.LeaseIDFromFlags(dcli.WithOwner(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			msg := &types.MsgCloseLease{
				LeaseID: id,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}
