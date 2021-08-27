package client

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ovrclk/akash/app"
	akash "github.com/ovrclk/akash/client"
)

// ReadPageRequest reads and builds the necessary page request flags for pagination.
func ReadPageRequest() (*query.PageRequest, error) {
	paginationOpts := PaginationOptsFromCmd()
	if paginationOpts.Page > 1 && paginationOpts.Offset > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "page and offset cannot be used together")
	}

	if paginationOpts.Page > 1 {
		paginationOpts.Offset = (paginationOpts.Page - 1) * paginationOpts.Limit
	}

	return &query.PageRequest{
		Key:        []byte(paginationOpts.PageKey),
		Offset:     paginationOpts.Offset,
		Limit:      paginationOpts.Limit,
		CountTotal: paginationOpts.CountTotal,
		Reverse:    paginationOpts.Reverse,
	}, nil
}

// GetClientQueryContext returns a Context from a command with fields set based on flags
// defined in AddQueryFlagsToCmd.
func GetClientQueryContext() (client.Context, error) {
	clientCtx := defaultContext()
	queryOpts := QueryFlagsFromCmd()
	if queryOpts.ChainID != "" {
		clientCtx = clientCtx.WithChainID(queryOpts.ChainID)
	}
	if queryOpts.Node != "" {
		clientCtx = clientCtx.WithNodeURI(queryOpts.Node)
		client, err := client.NewClientFromNode(queryOpts.Node)
		if err != nil {
			return clientCtx, err
		}
		clientCtx = clientCtx.WithClient(client)
	}
	if queryOpts.Height != 0 {
		clientCtx = clientCtx.WithHeight(queryOpts.Height)
	}
	if queryOpts.Output != "" {
		clientCtx = clientCtx.WithOutputFormat(queryOpts.Output)
	}

	return clientCtx, nil
}

func NewQueryClient(clientCtx client.Context) akash.QueryClient {
	return akash.NewQueryClientFromCtx(clientCtx)
}

func defaultContext() client.Context {
	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
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

	return clientCtx
}
