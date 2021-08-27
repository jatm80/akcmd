package deployment

import (
	"context"

	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/deployment/types"
	"github.com/ovrclk/akcmd/client"
	"github.com/ovrclk/akcmd/cmd/query/deployment/group"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "deployment",
		Desc: "Deployment query commands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{listCMD(), getCMD(), group.Cmd()},
	}

	return cmd
}

func listCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "list",
		Desc: "Query for all deployments",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddPaginationFlagsToCmd(cmd, "deployments")
			client.AddDeploymentFilterFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			dfilters := client.DepFiltersFromFlags()

			pageReq, err := client.ReadPageRequest()
			if err != nil {
				return err
			}

			params := &types.QueryDeploymentsRequest{
				Filters:    dfilters,
				Pagination: pageReq,
			}

			res, err := queryClient.Deployments(context.Background(), params)
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
		Desc: "Query deployment",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddDeploymentIDFlags(cmd)
			client.MarkReqDeploymentIDFlags(cmd)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext()
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			id := client.DeploymentIDFromFlags()

			res, err := queryClient.Deployment(context.Background(), &types.QueryDeploymentRequest{ID: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
