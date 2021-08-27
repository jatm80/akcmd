package order

import (
	"context"

	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/market/types"
	"github.com/ovrclk/akcmd/client"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "order",
		Desc: "Market order query commands",
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
		Desc: "Query for all orders",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddPaginationFlagsToCmd(cmd, "orders")
			client.AddOrderFilterFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			ofilters, err := client.OrderFiltersFromFlags()
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest()
			if err != nil {
				return err
			}

			params := &types.QueryOrdersRequest{
				Filters:    ofilters,
				Pagination: pageReq,
			}

			res, err := queryClient.Orders(context.Background(), params)
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
		Desc: "Query order",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddOrderIDFlags(cmd)
			client.MarkReqOrderIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			id, err := client.OrderIDFromFlags()
			if err != nil {
				return err
			}

			res, err := queryClient.Order(context.Background(), &types.QueryOrderRequest{ID: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Order)
		},
	}

	return cmd
}
