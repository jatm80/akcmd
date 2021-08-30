package client

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/gcli/v3"
	dcli "github.com/ovrclk/akash/x/deployment/client/cli"
	dtypes "github.com/ovrclk/akash/x/deployment/types"
	mtypes "github.com/ovrclk/akash/x/market/types"
)

var (
	qryOpts              = QueryOpts{}
	pageOpts             = PaginationOpts{}
	txOpts               = TransactionOpts{}
	deploymentFilterOpts = dtypes.DeploymentFilters{}
	deploymentIDOpts     = dtypes.DeploymentID{}
	gSeq                 uint64
	oSeq                 uint64
	provider             string
	deposit              string
	depositorAcc         string
	expiration           int64
)

type QueryOpts struct {
	ChainID string
	Node    string
	Height  int64
	Output  string
}

type PaginationOpts struct {
	Page       uint64
	PageKey    string
	Offset     uint64
	Limit      uint64
	CountTotal bool
	Reverse    bool
}

type TransactionOpts struct {
	ChainID          string
	Output           string
	KeyringDir       string
	From             string
	AccountNumber    uint64
	Sequence         uint64
	Note             string
	Fees             string
	GasPrices        string
	Node             string
	UseLedger        bool
	GasAdjustment    float64
	BroadcastMode    string
	DryRun           bool
	GenerateOnly     bool
	Offline          bool
	SkipConfirmation bool
	KeyringBackend   string
	SignMode         string
	TimeoutHeight    uint64
	FeeAccount       string
	Gas              string
}

func AddQueryFlagsToCmd(cmd *gcli.Command) {
	cmd.StrOpt(&qryOpts.ChainID, flags.FlagChainID, "", "", "The network chain ID")
	cmd.StrOpt(&qryOpts.Node, flags.FlagNode, "", "tcp://localhost:26657",
		"<host>:<port> to Tendermint RPC interface for this chain")
	cmd.Int64Opt(&qryOpts.Height, flags.FlagHeight, "", 0,
		"Use a specific height to query state at (this can error if the node is pruning state)")
	cmd.StrOpt(&qryOpts.Output, "output", "o", "text", "Output format (text|json)")

	cmd.Required(flags.FlagChainID)
}

func QueryFlagsFromCmd() QueryOpts {
	return qryOpts
}

func AddTxFlagsToCmd(cmd *gcli.Command) {
	cmd.StrOpt(&txOpts.ChainID, flags.FlagChainID, "", "", "The network chain ID")
	cmd.StrOpt(&txOpts.Output, "output", "o", "json", "Output format (text|json)")
	cmd.StrOpt(&txOpts.KeyringDir, flags.FlagKeyringDir, "", "",
		"The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.StrOpt(&txOpts.From, flags.FlagFrom, "", "",
		"Name or address of private key with which to sign")
	cmd.Uint64Opt(&txOpts.AccountNumber, flags.FlagAccountNumber, "a", 0,
		"The account number of the signing account (offline mode only)")
	cmd.Uint64Opt(&txOpts.Sequence, flags.FlagSequence, "s", 0,
		"The sequence number of the signing account (offline mode only)")
	cmd.StrOpt(&txOpts.Note, flags.FlagNote, "", "",
		"Note to add a description to the transaction (previously --memo)")
	cmd.StrOpt(&txOpts.Fees, flags.FlagFees, "", "",
		"Fees to pay along with transaction; eg: 10uatom")
	cmd.StrOpt(&txOpts.GasPrices, flags.FlagGasPrices, "", "",
		"Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)")
	cmd.StrOpt(&txOpts.Node, flags.FlagNode, "", "tcp://localhost:26657",
		"<host>:<port> to tendermint rpc interface for this chain")
	cmd.BoolOpt(&txOpts.UseLedger, flags.FlagUseLedger, "", false, "Use a connected Ledger device")
	cmd.Float64Opt(&txOpts.GasAdjustment, flags.FlagGasAdjustment, "", flags.DefaultGasAdjustment,
		"adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored ")
	cmd.StrOpt(&txOpts.BroadcastMode, flags.FlagBroadcastMode, "b", flags.BroadcastSync,
		"Transaction broadcasting mode (sync|async|block)")
	cmd.BoolOpt(&txOpts.DryRun, flags.FlagDryRun, "", false,
		"ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it")
	cmd.BoolOpt(&txOpts.GenerateOnly, flags.FlagGenerateOnly, "", false,
		"Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible)")
	cmd.BoolOpt(&txOpts.Offline, flags.FlagOffline, "", false,
		"Offline mode (does not allow any online functionality")
	cmd.BoolOpt(&txOpts.SkipConfirmation, flags.FlagSkipConfirmation, "y", false,
		"Skip tx broadcasting prompt confirmation")
	cmd.StrOpt(&txOpts.KeyringBackend, flags.FlagKeyringBackend, "", flags.DefaultKeyringBackend,
		"Select keyring's backend (os|file|kwallet|pass|test|memory)")
	cmd.StrOpt(&txOpts.SignMode, flags.FlagSignMode, "", "",
		"Choose sign mode (direct|amino-json), this is an advanced feature")
	cmd.Uint64Opt(&txOpts.TimeoutHeight, flags.FlagTimeoutHeight, "", 0,
		"Set a block timeout height to prevent the tx from being committed past a certain height")
	cmd.StrOpt(&txOpts.FeeAccount, flags.FlagFeeAccount, "", "",
		"Fee account pays fees for the transaction instead of deducting from the signer")

	// --gas can accept integers and "auto"
	cmd.StrOpt(&txOpts.Gas, flags.FlagGas, "", "",
		fmt.Sprintf("gas limit to set per-transaction; set to %q to calculate sufficient gas automatically (default %d)", flags.GasFlagAuto, flags.DefaultGasLimit))

	cmd.Required(flags.FlagChainID)
}

func TxFlagsFromCmd() TransactionOpts {
	return txOpts
}

func AddPaginationFlagsToCmd(cmd *gcli.Command, query string) {
	cmd.Uint64Opt(&pageOpts.Page, flags.FlagPage, "", 1,
		fmt.Sprintf("pagination page of %s to query for. This sets offset to a multiple of limit", query))
	cmd.StrOpt(&pageOpts.PageKey, flags.FlagPageKey, "", "",
		fmt.Sprintf("pagination page-key of %s to query for", query))
	cmd.Uint64Opt(&pageOpts.Offset, flags.FlagOffset, "", 0,
		fmt.Sprintf("pagination offset of %s to query for", query))
	cmd.Uint64Opt(&pageOpts.Limit, flags.FlagLimit, "", 100,
		fmt.Sprintf("pagination limit of %s to query for", query))
	cmd.BoolOpt(&pageOpts.CountTotal, flags.FlagCountTotal, "", false,
		fmt.Sprintf("count total number of records in %s to query for", query))
	cmd.BoolOpt(&pageOpts.Reverse, flags.FlagReverse, "", false,
		"results are sorted in descending order")
}

func PaginationOptsFromCmd() PaginationOpts {
	return pageOpts
}

func AddDeploymentFilterFlags(cmd *gcli.Command) {
	cmd.StrOpt(&deploymentFilterOpts.Owner, "owner", "", "", "deployment owner address to filter")
	cmd.StrOpt(&deploymentFilterOpts.State, "state", "", "", "deployment state to filter (active,closed)")
	cmd.Uint64Opt(&deploymentFilterOpts.DSeq, "dseq", "", 0, "deployment sequence to filter")
}

func DepFiltersFromFlags() dtypes.DeploymentFilters {
	return deploymentFilterOpts
}

func AddDeploymentIDFlags(cmd *gcli.Command) {
	cmd.StrOpt(&deploymentIDOpts.Owner, "owner", "", "", "Deployment Owner")
	cmd.Uint64Opt(&deploymentIDOpts.DSeq, "dseq", "", 0, "Deployment Sequence")
}

func MarkReqDeploymentIDFlags(cmd *gcli.Command, opts ...dcli.DeploymentIDOption) {
	// TODO: update the code here once dcli.deploymentIDOption{} is made public in akash repo
	cmd.Required("owner", "dseq")
}

func DeploymentIDFromFlags(opts ...dcli.MarketOption) (dtypes.DeploymentID, error) {
	opt := &dcli.MarketOptions{}

	for _, o := range opts {
		o(opt)
	}

	// if --owner flag was explicitly provided, use that.
	if deploymentIDOpts.Owner != "" {
		var err error
		opt.Owner, err = sdk.AccAddressFromBech32(deploymentIDOpts.Owner)
		if err != nil {
			return deploymentIDOpts, err
		}
	}
	deploymentIDOpts.Owner = opt.Owner.String()

	return deploymentIDOpts, nil
}

func AddGroupIDFlags(cmd *gcli.Command) {
	AddDeploymentIDFlags(cmd)
	cmd.Uint64Opt(&gSeq, "gseq", "", 1, "Group Sequence")
}

func MarkReqGroupIDFlags(cmd *gcli.Command, opts ...dcli.DeploymentIDOption) {
	MarkReqDeploymentIDFlags(cmd, opts...)
}

func GroupIDFromFlags(opts ...dcli.MarketOption) (dtypes.GroupID, error) {
	dID, err := DeploymentIDFromFlags(opts...)
	if err != nil {
		return dtypes.GroupID{}, err
	}

	val, err := getGSeq()
	if err != nil {
		return dtypes.GroupID{}, err
	}

	return dtypes.MakeGroupID(dID, val), nil
}

func getGSeq() (uint32, error) {
	if gSeq > math.MaxUint32 {
		return 0, errors.New("gseq out of uint32 range")
	}
	return uint32(gSeq), nil
}

func AddOrderFilterFlags(cmd *gcli.Command) {
	cmd.StrOpt(&deploymentFilterOpts.Owner, "owner", "", "", "order owner address to filter")
	cmd.StrOpt(&deploymentFilterOpts.State, "state", "", "", "order state to filter (open,matched,closed)")
	cmd.Uint64Opt(&deploymentFilterOpts.DSeq, "dseq", "", 0, "deployment sequence to filter")
	cmd.Uint64Opt(&gSeq, "gseq", "", 1, "group sequence to filter")
	cmd.Uint64Opt(&oSeq, "oseq", "", 1, "order sequence to filter")
}

func OrderFiltersFromFlags() (mtypes.OrderFilters, error) {
	dFilters := DepFiltersFromFlags()
	filter := mtypes.OrderFilters{
		Owner: dFilters.Owner,
		State: dFilters.State,
		DSeq:  dFilters.DSeq,
	}
	var err error
	if filter.GSeq, err = getGSeq(); err != nil {
		return filter, err
	}
	if filter.OSeq, err = getOSeq(); err != nil {
		return filter, err
	}
	return filter, nil
}

func getOSeq() (uint32, error) {
	if oSeq > math.MaxUint32 {
		return 0, errors.New("oseq out of uint32 range")
	}
	return uint32(oSeq), nil
}

func AddOrderIDFlags(cmd *gcli.Command) {
	AddGroupIDFlags(cmd)
	cmd.Uint64Opt(&oSeq, "oseq", "", 1, "Order Sequence")
}

func MarkReqOrderIDFlags(cmd *gcli.Command, opts ...dcli.DeploymentIDOption) {
	MarkReqGroupIDFlags(cmd, opts...)
}

func OrderIDFromFlags(opts ...dcli.MarketOption) (mtypes.OrderID, error) {
	gID, err := GroupIDFromFlags(opts...)
	if err != nil {
		return mtypes.OrderID{}, err
	}
	val, err := getOSeq()
	if err != nil {
		return mtypes.OrderID{}, err
	}
	return mtypes.MakeOrderID(gID, val), nil
}

func AddBidFilterFlags(cmd *gcli.Command) {
	cmd.StrOpt(&deploymentFilterOpts.Owner, "owner", "", "", "bid owner address to filter")
	cmd.StrOpt(&deploymentFilterOpts.State, "state", "", "", "bid state to filter (open,matched,lost,closed)")
	cmd.Uint64Opt(&deploymentFilterOpts.DSeq, "dseq", "", 0, "deployment sequence to filter")
	cmd.Uint64Opt(&gSeq, "gseq", "", 1, "group sequence to filter")
	cmd.Uint64Opt(&oSeq, "oseq", "", 1, "order sequence to filter")
	cmd.StrOpt(&provider, "provider", "", "", "bid provider address to filter")
}

func BidFiltersFromFlags() (mtypes.BidFilters, error) {
	oFilters, err := OrderFiltersFromFlags()
	if err != nil {
		return mtypes.BidFilters{}, err
	}
	bFilters := mtypes.BidFilters{
		Owner: oFilters.Owner,
		DSeq:  oFilters.DSeq,
		GSeq:  oFilters.GSeq,
		OSeq:  oFilters.OSeq,
		State: oFilters.State,
	}
	bFilters.Provider, err = getProviderFilter()
	if err != nil {
		return mtypes.BidFilters{}, err
	}
	return bFilters, nil
}

func getProviderFilter() (string, error) {
	if provider != "" {
		_, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return "", err
		}
	}
	return provider, nil
}

func AddProviderFlag(cmd *gcli.Command) {
	cmd.StrOpt(&provider, "provider", "", "", "Provider")
}

func MarkReqProviderFlag(cmd *gcli.Command) {
	cmd.Required("provider")
}

func ProviderFromFlag() (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32(provider)
}

func AddBidIDFlags(cmd *gcli.Command) {
	AddOrderIDFlags(cmd)
	AddProviderFlag(cmd)
}

func AddQueryBidIDFlags(cmd *gcli.Command) {
	AddBidIDFlags(cmd)
}

func MarkReqBidIDFlags(cmd *gcli.Command, opts ...dcli.DeploymentIDOption) {
	MarkReqOrderIDFlags(cmd, opts...)
	MarkReqProviderFlag(cmd)
}

func BidIDFromFlags(opts ...dcli.MarketOption) (mtypes.BidID, error) {
	prev, err := OrderIDFromFlags(opts...)
	if err != nil {
		return mtypes.BidID{}, err
	}

	opt := &dcli.MarketOptions{}

	for _, o := range opts {
		o(opt)
	}

	if opt.Provider.Empty() {
		if opt.Provider, err = ProviderFromFlag(); err != nil {
			return mtypes.BidID{}, err
		}
	}

	return mtypes.MakeBidID(prev, opt.Provider), nil
}

func AddLeaseIDFlags(cmd *gcli.Command) {
	AddBidIDFlags(cmd)
}
func MarkReqLeaseIDFlags(cmd *gcli.Command, opts ...dcli.DeploymentIDOption) {
	MarkReqBidIDFlags(cmd, opts...)
}

func LeaseIDFromFlags(opts ...dcli.MarketOption) (mtypes.LeaseID, error) {
	bid, err := BidIDFromFlags(opts...)
	if err != nil {
		return mtypes.LeaseID{}, err
	}

	return bid.LeaseID(), nil
}

func AddLeaseFilterFlags(cmd *gcli.Command) {
	cmd.StrOpt(&deploymentFilterOpts.Owner, "owner", "", "", "lease owner address to filter")
	cmd.StrOpt(&deploymentFilterOpts.State, "state", "", "", "lease state to filter (active,insufficient_funds,closed)")
	cmd.Uint64Opt(&deploymentFilterOpts.DSeq, "dseq", "", 0, "deployment sequence to filter")
	cmd.Uint64Opt(&gSeq, "gseq", "", 1, "group sequence to filter")
	cmd.Uint64Opt(&oSeq, "oseq", "", 1, "order sequence to filter")
	cmd.StrOpt(&provider, "provider", "", "", "bid provider address to filter")
}

func LeaseFiltersFromFlags() (mtypes.LeaseFilters, error) {
	bFilters, err := BidFiltersFromFlags()
	if err != nil {
		return mtypes.LeaseFilters{}, err
	}
	return mtypes.LeaseFilters(bFilters), nil
}

func AddDepositFlags(cmd *gcli.Command, dflt sdk.Coin) {
	cmd.StrOpt(&deposit, "deposit", "", dflt.String(), "Deposit amount")
}

func DepositFromFlags() (sdk.Coin, error) {
	return sdk.ParseCoinNormalized(deposit)
}

// AddDepositorFlag adds the `--depositor-account` flag
func AddDepositorFlag(cmd *gcli.Command) {
	cmd.StrOpt(&depositorAcc, "depositor-account", "", "",
		"Depositor account pays for the deposit instead of deducting from the owner")
}

// DepositorFromFlags returns the depositor account if one was specified in flags,
// otherwise it returns the owner's account.
func DepositorFromFlags(owner string) (string, error) {
	// if no depositor is specified, owner is the default depositor
	if strings.TrimSpace(depositorAcc) == "" {
		return owner, nil
	}

	_, err := sdk.AccAddressFromBech32(depositorAcc)
	return depositorAcc, err
}

func AddExpirationFlag(cmd *gcli.Command) {
	cmd.Int64Opt(&expiration, "expiration", "", time.Now().AddDate(1, 0, 0).Unix(),
		"The Unix timestamp. Default is one year.")
}

func ExpirationFromFlag() int64 {
	return expiration
}
