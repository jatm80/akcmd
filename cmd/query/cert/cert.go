package cert

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/cert/types"
	"github.com/ovrclk/akcmd/client"
	"github.com/pkg/errors"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "cert",
		Desc: "Certificate query commands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{listCMD()},
	}

	return cmd
}

var listOpts = struct {
	Serial string
	Owner  string
	State  string
}{}

func listCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "list",
		Desc: "Query for all certificates",
		Config: func(cmd *gcli.Command) {
			client.AddQueryFlagsToCmd(cmd)
			client.AddPaginationFlagsToCmd(cmd, "certificates")
			cmd.StrOpt(&listOpts.Serial, "serial", "", "", "filter certificates by serial number")
			cmd.StrOpt(&listOpts.Owner, "owner", "", "", "filter certificates by owner")
			cmd.StrOpt(&listOpts.State, "state", "", "", "filter certificates by valid|revoked")
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

			params := &types.QueryCertificatesRequest{
				Pagination: pageReq,
			}

			if listOpts.Owner != "" {
				var owner sdk.Address
				if owner, err = sdk.AccAddressFromBech32(listOpts.Owner); err != nil {
					return err
				}

				params.Filter.Owner = owner.String()
			}

			if listOpts.Serial != "" {
				if params.Filter.Owner == "" {
					return errors.Errorf("--serial flag requires --owner to be set")
				}
				val, valid := new(big.Int).SetString(listOpts.Serial, 10)
				if !valid {
					return errors.New("invalid value in serial flag. expected integer")
				}

				params.Filter.Serial = val.String()
			}

			if listOpts.State != "" {
				if listOpts.State != "valid" && listOpts.State != "revoked" {
					return errors.Errorf("invalid value of --state flag. expected valid|revoked")
				}

				params.Filter.State = listOpts.State
			}

			res, err := queryClient.Certificates(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
