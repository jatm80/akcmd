package group

import (
	"context"

	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/deployment/types"
	"github.com/ovrclk/akcmd/client"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "group",
		Desc: "Deployment group query commands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{getCMD()},
	}

	return cmd
}

func getCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "get",
		Desc: "Query group of deployment",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddGroupIDFlags(cmd)
			client.MarkReqGroupIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			id, err := client.GroupIDFromFlags()
			if err != nil {
				return err
			}

			res, err := queryClient.Group(context.Background(), &types.QueryGroupRequest{ID: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Group)
		},
	}

	return cmd
}
