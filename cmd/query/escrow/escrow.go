package escrow

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/deployment/client/cli"
	deploymentTypes "github.com/ovrclk/akash/x/deployment/types"
	marketTypes "github.com/ovrclk/akash/x/market/types"
	"github.com/ovrclk/akcmd/client"
	"gopkg.in/yaml.v2"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "escrow",
		Desc: "Escrow query commands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{blocksRemainingCMD()},
	}

	return cmd
}

// Define 6.5 seconds as average block-time
const secondsPerBlock = 6.5

var errNoLeaseMatches = errors.New("leases for deployment do not exist")

func blocksRemainingCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "blocks-remaining",
		Desc: "Compute the number of blocks remaining for an escrow account",
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

			marketClient := marketTypes.NewQueryClient(clientCtx)
			ctx := context.Background()

			id := client.DeploymentIDFromFlags()

			// Fetch leases matching owner & dseq
			leaseRequest := marketTypes.QueryLeasesRequest{
				Filters: marketTypes.LeaseFilters{
					Owner:    id.Owner,
					DSeq:     id.DSeq,
					GSeq:     0,
					OSeq:     0,
					Provider: "",
					State:    "active",
				},
				Pagination: nil,
			}

			leasesResponse, err := marketClient.Leases(ctx, &leaseRequest)
			if err != nil {
				return err
			}

			leases := make([]marketTypes.Lease, 0)
			for _, lease := range leasesResponse.Leases {
				leases = append(leases, lease.Lease)
			}

			// Fetch the balance of the escrow account
			deploymentClient := deploymentTypes.NewQueryClient(clientCtx)
			totalLeaseAmount := cosmosTypes.NewInt(0)
			blockchainHeight, err := cli.CurrentBlockHeight(clientCtx)
			if err != nil {
				return err
			}
			if 0 == len(leases) {
				return errNoLeaseMatches
			}
			for _, lease := range leases {

				amount := lease.Price.Amount
				totalLeaseAmount = totalLeaseAmount.Add(amount)

			}
			res, err := deploymentClient.Deployment(ctx, &deploymentTypes.QueryDeploymentRequest{
				ID: deploymentTypes.DeploymentID{Owner: id.Owner, DSeq: id.DSeq},
			})
			if err != nil {
				return err
			}
			balance := res.EscrowAccount.TotalBalance().Amount
			settledAt := res.EscrowAccount.SettledAt
			balanceRemain := float64(balance.Int64() - ((int64(blockchainHeight) - settledAt) * (totalLeaseAmount.Int64())))
			blocksRemain := (balanceRemain / float64(totalLeaseAmount.Int64()))
			output := struct {
				BalanceRemain       float64 `json:"balance_remaining" yaml:"balance_remaining"`
				BlocksRemain        float64 `json:"blocks_remaining" yaml:"blocks_remaining"`
				EstimatedTimeRemain string  `json:"estimated_time_remaining" yaml:"estimated_time_remaining"`
			}{
				BalanceRemain:       balanceRemain,
				BlocksRemain:        blocksRemain,
				EstimatedTimeRemain: (time.Duration(math.Floor(secondsPerBlock*blocksRemain)) * time.Second).String(),
			}

			outputType := client.QueryFlagsFromCmd().Output

			var data []byte
			if outputType == "json" {
				data, err = json.MarshalIndent(output, " ", "\t")
			} else {
				data, err = yaml.Marshal(output)
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintBytes(data)
		},
	}

	return cmd
}
