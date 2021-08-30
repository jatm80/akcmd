package group

import (
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/deployment/types"
	"github.com/ovrclk/akcmd/client"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "group",
		Desc: "Modify a Deployment's specific Group",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{closeCMD(), pauseCMD(), startCMD()},
	}

	return cmd
}

func closeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "close",
		Desc: "close a Deployment's specific Group",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddGroupIDFlags(cmd)
			client.MarkReqGroupIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.GroupIDFromFlags()
			if err != nil {
				return err
			}

			msg := &types.MsgCloseGroup{
				ID: id,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func pauseCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "pause",
		Desc: "pause a Deployment's specific Group",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddGroupIDFlags(cmd)
			client.MarkReqGroupIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.GroupIDFromFlags()
			if err != nil {
				return err
			}

			msg := &types.MsgPauseGroup{
				ID: id,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func startCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "start",
		Desc: "start a Deployment's specific Group",
		Config: func(cmd *gcli.Command) {
			client.AddTxFlagsToCmd(cmd)
			client.AddGroupIDFlags(cmd)
			client.MarkReqGroupIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			id, err := client.GroupIDFromFlags()
			if err != nil {
				return err
			}

			msg := &types.MsgStartGroup{
				ID: id,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}
