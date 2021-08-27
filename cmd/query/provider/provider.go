package provider

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/provider/types"
	"github.com/ovrclk/akcmd/client"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "provider",
		Desc: "Provider query commands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{listCMD(), getCMD()},
	}

	return cmd
}

func listCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "list",
		Desc: "Query for all providers",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddPaginationFlagsToCmd(cmd, "providers")
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest()
			if err != nil {
				return err
			}

			params := &types.QueryProvidersRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Providers(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func getCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "get",
		Desc: "Query provider",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)

			cmd.AddArg("owner", "owner address", true)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Provider(context.Background(), &types.QueryProviderRequest{Owner: owner.String()})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Provider)
		},
	}

	return cmd
}
