package client

import (
	"context"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/app"
	akash "github.com/ovrclk/akash/client"
	"github.com/spf13/cobra"
)

func DefaultContext() client.Context {
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		//WithHomeDir(app.DefaultHome).
		WithOffline(false)
	//WithOutputFormat("json").
	//WithKeyringDir().
	//WithChainID().
	//WithKeyring().
	//WithNodeURI()
	//WithClient()

	return initClientCtx
}

func QueryClient() akash.QueryClient {
	cmd := withDefaultContext(&cobra.Command{})
	cctx := client.GetClientContextFromCmd(cmd)
	qc := akash.NewQueryClientFromCtx(cctx)
	return qc
}

func ReadPageRequest(cmd *gcli.Command) (*query.PageRequest, error) {
	return client.ReadPageRequest(toCobraCmd(cmd).Flags())
}

func GetClientQueryContext(cmd *gcli.Command) (client.Context, error) {
	return client.GetClientQueryContext(withDefaultContext(toCobraCmd(cmd)))
}

func NewQueryClient(ctx client.Context) akash.QueryClient {
	return akash.NewQueryClientFromCtx(ctx)
}

func toCobraCmd(cmd *gcli.Command) *cobra.Command {
	c := &cobra.Command{Use: cmd.Name}
	c.Flags().AddGoFlagSet(cmd.FSet())
	return c
}

func withDefaultContext(cmd *cobra.Command) *cobra.Command {
	ctx := context.WithValue(context.Background(), client.ClientContextKey, &client.Context{})
	if err := cmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}
	if err := client.SetCmdClientContextHandler(DefaultContext(), cmd); err != nil {
		panic(err)
	}
	return cmd
}
