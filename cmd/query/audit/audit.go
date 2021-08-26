package audit

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/audit/types"
	"github.com/ovrclk/akcmd/client"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "audit",
		Desc: "Audit query commands",
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
		Desc: "Query for all provider attributes",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddPaginationFlagsToCmd(cmd, "providers")
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd)
			if err != nil {
				return err
			}

			params := &types.QueryAllProvidersAttributesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllProvidersAttributes(context.Background(), params)
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
		Desc: "Query provider attributes",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)

			cmd.AddArg("owner", "owner address", true)
			cmd.AddArg("auditor", "auditor address", false)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var res *types.QueryProvidersResponse
			if len(args) == 1 {
				res, err = queryClient.ProviderAttributes(context.Background(),
					&types.QueryProviderAttributesRequest{
						Owner: owner.String(),
					},
				)
			} else {
				var auditor sdk.AccAddress
				if auditor, err = sdk.AccAddressFromBech32(args[1]); err != nil {
					return err
				}

				res, err = queryClient.ProviderAuditorAttributes(context.Background(),
					&types.QueryProviderAuditorRequest{
						Auditor: auditor.String(),
						Owner:   owner.String(),
					},
				)
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
