package attr

import (
	"context"
	"errors"
	"sort"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/gcli/v3"
	akashtypes "github.com/ovrclk/akash/types"
	"github.com/ovrclk/akash/x/audit/types"
	ptypes "github.com/ovrclk/akash/x/provider/types"
	"github.com/ovrclk/akcmd/client"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "attr",
		Desc: "Manage provider attributes",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{createCMD(), deleteCMD()},
	}

	return cmd
}

func createCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name:   "create",
		Desc:   "Create/update provider attributes",
		Config: config,
		Func: func(cmd *gcli.Command, args []string) error {
			if ((len(args) - 1) % 2) != 0 {
				return errors.New("attributes must be provided as pairs")
			}

			providerAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			attr, err := readAttributes(clientCtx, providerAddress.String(), args[1:])
			if err != nil {
				return err
			}

			if len(attr) == 0 {
				return errors.New("no attributes provided|found")
			}

			msg := &types.MsgSignProviderAttributes{
				Auditor:    clientCtx.GetFromAddress().String(),
				Owner:      providerAddress.String(),
				Attributes: attr,
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func deleteCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name:   "delete",
		Desc:   "Delete provider attributes",
		Config: config,
		Func: func(cmd *gcli.Command, args []string) error {
			providerAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			keys, err := readKeys(args[1:])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			msg := &types.MsgDeleteProviderAttributes{
				Auditor: clientCtx.GetFromAddress().String(),
				Owner:   providerAddress.String(),
				Keys:    keys,
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return client.BroadcastTX(clientCtx, msg)
		},
	}

	return cmd
}

func config(cmd *gcli.Command) {
	client.AddTxFlagsToCmd(cmd)
	cmd.Required(flags.FlagFrom)

	cmd.AddArg("provider", "Provider address", true)
	cmd.AddArg("attributes", "Provider attributes", false, true)
}

// readAttributes try read attributes from both cobra arguments or query
// if no arguments were provided then query provider and sign all found
// read from stdin uses trick to check if it's file descriptor is a pipe
// which happens when some data is piped for example cat attr.yaml | akash ...
func readAttributes(cctx sdkclient.Context, provider string, args []string) (akashtypes.Attributes, error) {
	var attr akashtypes.Attributes

	if len(args) != 0 {
		for i := 0; i < len(args); i += 2 {
			attr = append(attr, akashtypes.Attribute{
				Key:   args[i],
				Value: args[i+1],
			})
		}
	} else {
		resp, err := ptypes.NewQueryClient(cctx).Provider(context.Background(), &ptypes.QueryProviderRequest{Owner: provider})
		if err != nil {
			return nil, err
		}

		attr = append(attr, resp.Provider.Attributes...)
	}

	sort.SliceStable(attr, func(i, j int) bool {
		return attr[i].Key < attr[j].Value
	})

	if checkAttributeDuplicates(attr) {
		return nil, errors.New("supplied attributes with duplicate keys")
	}

	return attr, nil
}

func readKeys(args []string) ([]string, error) {
	sort.SliceStable(args, func(i, j int) bool {
		return args[i] < args[j]
	})

	if checkKeysDuplicates(args) {
		return nil, errors.New("supplied attributes with duplicate keys")
	}

	return args, nil
}

func checkAttributeDuplicates(attr akashtypes.Attributes) bool {
	keys := make(map[string]bool)

	for _, entry := range attr {
		if _, value := keys[entry.Key]; !value {
			keys[entry.Key] = true
		} else {
			return true
		}
	}
	return false
}

func checkKeysDuplicates(k []string) bool {
	keys := make(map[string]bool)

	for _, entry := range k {
		if _, value := keys[entry]; !value {
			keys[entry] = true
		} else {
			return true
		}
	}
	return false
}
