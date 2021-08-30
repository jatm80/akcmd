package bid

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/gcli/v3"
	dcli "github.com/ovrclk/akash/x/deployment/client/cli"
	"github.com/ovrclk/akash/x/market/client/cli"
	"github.com/ovrclk/akash/x/market/types"
	"github.com/ovrclk/akcmd/client"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "bid",
		Desc: "Bid subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{createCMD(), closeCMD()},
	}

	return cmd
}

var createOpts = struct {
	price string
}{}

func createCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "Create a market bid",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddOrderIDFlags(cmd)
			cmd.StrOpt(&createOpts.price, "price", "", "", "Bid Price")
			client.AddDepositFlags(cmd, cli.DefaultDeposit)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinNormalized(createOpts.price)
			if err != nil {
				return err
			}

			id, err := client.OrderIDFromFlags(dcli.WithProvider(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			deposit, err := client.DepositFromFlags()
			if err != nil {
				return err
			}

			msg := &types.MsgCreateBid{
				Order:    id,
				Provider: clientCtx.GetFromAddress().String(),
				Price:    coins,
				Deposit:  deposit,
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
		Desc: "Close a market bid",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddBidIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.BidIDFromFlags(dcli.WithProvider(clientCtx.FromAddress))
			if err != nil {
				return err
			}

			msg := &types.MsgCloseBid{
				BidID: id,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}
