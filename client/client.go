package client

import (
	"context"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ovrclk/akash/app"
	akash "github.com/ovrclk/akash/client"
	"github.com/ovrclk/akash/client/broadcaster"
)

// ReadPageRequest reads and builds the necessary page request flags for pagination.
func ReadPageRequest() (*query.PageRequest, error) {
	opts := PaginationOptsFromCmd()
	if opts.Page > 1 && opts.Offset > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "page and offset cannot be used together")
	}

	if opts.Page > 1 {
		opts.Offset = (opts.Page - 1) * opts.Limit
	}

	return &query.PageRequest{
		Key:        []byte(opts.PageKey),
		Offset:     opts.Offset,
		Limit:      opts.Limit,
		CountTotal: opts.CountTotal,
		Reverse:    opts.Reverse,
	}, nil
}

// GetClientQueryContext returns a Context from a command with fields set based on flags
// defined in AddQueryFlagsToCmd.
func GetClientQueryContext() (client.Context, error) {
	clientCtx := defaultContext()
	opts := QueryFlagsFromCmd()
	if opts.ChainID != "" {
		clientCtx = clientCtx.WithChainID(opts.ChainID)
	}
	if opts.Node != "" {
		clientCtx = clientCtx.WithNodeURI(opts.Node)
		client, err := client.NewClientFromNode(opts.Node)
		if err != nil {
			return clientCtx, err
		}
		clientCtx = clientCtx.WithClient(client)
	}
	if opts.Height != 0 {
		clientCtx = clientCtx.WithHeight(opts.Height)
	}
	if opts.Output != "" {
		clientCtx = clientCtx.WithOutputFormat(opts.Output)
	}

	return clientCtx, nil
}

// GetClientTxContext returns a Context from a command with fields set based on flags
// defined in AddTxFlagsToCmd.
func GetClientTxContext() (client.Context, error) {
	clientCtx := defaultContext()
	opts := TxFlagsFromCmd()

	if opts.Output != "" {
		clientCtx = clientCtx.WithOutputFormat(opts.Output)
	}

	clientCtx = clientCtx.WithSimulation(opts.DryRun)

	if opts.KeyringDir == "" {
		opts.KeyringDir = clientCtx.HomeDir
	}
	clientCtx = clientCtx.WithKeyringDir(opts.KeyringDir)

	if opts.ChainID != "" {
		clientCtx = clientCtx.WithChainID(opts.ChainID)
	}

	if opts.KeyringBackend != "" {
		kr, err := client.NewKeyringFromBackend(clientCtx, opts.KeyringBackend)
		if err != nil {
			return clientCtx, err
		}

		clientCtx = clientCtx.WithKeyring(kr)
	}

	if opts.Node != "" {
		clientCtx = clientCtx.WithNodeURI(opts.Node)
		client, err := client.NewClientFromNode(opts.Node)
		if err != nil {
			return clientCtx, err
		}
		clientCtx = clientCtx.WithClient(client)
	}

	clientCtx = clientCtx.WithGenerateOnly(opts.GenerateOnly)
	clientCtx = clientCtx.WithOffline(opts.Offline)
	clientCtx = clientCtx.WithUseLedger(opts.UseLedger)
	clientCtx = clientCtx.WithBroadcastMode(opts.BroadcastMode)
	clientCtx = clientCtx.WithSkipConfirmation(opts.SkipConfirmation)
	clientCtx = clientCtx.WithSignModeStr(opts.SignMode)

	if opts.FeeAccount != "" {
		granterAcc, err := sdk.AccAddressFromBech32(opts.FeeAccount)
		if err != nil {
			return clientCtx, err
		}

		clientCtx = clientCtx.WithFeeGranterAddress(granterAcc)
	}

	{
		fromAddr, fromName, keyType, err := client.GetFromFields(clientCtx.Keyring, opts.From, clientCtx.GenerateOnly)
		if err != nil {
			return clientCtx, err
		}
		clientCtx = clientCtx.WithFrom(opts.From).WithFromAddress(fromAddr).WithFromName(fromName)

		// If the `from` signer account is a ledger key, we need to use
		// SIGN_MODE_AMINO_JSON, because ledger doesn't support proto yet.
		// ref: https://github.com/cosmos/cosmos-sdk/issues/8109
		if keyType == keyring.TypeLedger && clientCtx.SignModeStr != flags.SignModeLegacyAminoJSON {
			fmt.Println("Default sign-mode 'direct' not supported by Ledger, using sign-mode 'amino-json'.")
			clientCtx = clientCtx.WithSignModeStr(flags.SignModeLegacyAminoJSON)
		}
	}

	return clientCtx, nil
}

func GetServerContext() *server.Context {
	return server.NewDefaultContext()
}

func NewQueryClient(clientCtx client.Context) akash.QueryClient {
	return akash.NewQueryClientFromCtx(clientCtx)
}

func BroadcastTX(clientCtx client.Context, msgs ...sdk.Msg) error {
	broadcaster, err := newTxBroadcaster(clientCtx)
	if err != nil {
		return nil
	}

	return broadcaster.Broadcast(context.Background(), msgs...)
}

func newTxBroadcaster(clientCtx client.Context) (broadcaster.Client, error) {
	txFactory := newFactoryCLI(clientCtx)
	info, err := txFactory.Keybase().Key(clientCtx.GetFromName())
	if err != nil {
		return nil, err
	}
	return broadcaster.NewClient(clientCtx, txFactory, info), nil
}

// newFactoryCLI creates a new Factory.
func newFactoryCLI(clientCtx client.Context) tx.Factory {
	signModeStr := clientCtx.SignModeStr

	signMode := signing.SignMode_SIGN_MODE_UNSPECIFIED
	switch signModeStr {
	case flags.SignModeDirect:
		signMode = signing.SignMode_SIGN_MODE_DIRECT
	case flags.SignModeLegacyAminoJSON:
		signMode = signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	}

	opts := TxFlagsFromCmd()
	gasSetting, _ := flags.ParseGasSetting(opts.Gas)

	f := tx.Factory{}.
		WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever).
		WithKeybase(clientCtx.Keyring).WithChainID(clientCtx.ChainID).WithGas(gasSetting.Gas).
		WithSimulateAndExecute(gasSetting.Simulate).WithAccountNumber(opts.AccountNumber).
		WithSequence(opts.Sequence).WithTimeoutHeight(opts.TimeoutHeight).
		WithGasAdjustment(opts.GasAdjustment).WithMemo(opts.Note).WithSignMode(signMode).
		WithFees(opts.Fees).WithGasPrices(opts.GasPrices)

	return f
}

func defaultContext() client.Context {
	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
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
