package lease

import (
	"context"

	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/market/types"
	"github.com/ovrclk/akcmd/client"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "lease",
		Desc: "Market lease query commands",
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
		Desc: "Query for all leases",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddPaginationFlagsToCmd(cmd, "leases")
			client.AddLeaseFilterFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			lfilters, err := client.LeaseFiltersFromFlags()
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest()
			if err != nil {
				return err
			}

			params := &types.QueryLeasesRequest{
				Filters:    lfilters,
				Pagination: pageReq,
			}

			res, err := queryClient.Leases(context.Background(), params)
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
		Desc: "Query lease",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddQueryBidIDFlags(cmd)
			client.MarkReqBidIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			bidID, err := client.BidIDFromFlags()
			if err != nil {
				return err
			}

			res, err := queryClient.Lease(context.Background(), &types.QueryLeaseRequest{ID: types.MakeLeaseID(bidID)})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
